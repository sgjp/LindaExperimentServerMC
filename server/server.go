package server

import (//"reflect"
	"github.com/sgjp/LindaExperimentServerMC/tupleSpace"
	"log"
	"github.com/sgjp/go-coap"
	"strings"

	//"time"
	"net"
	"strconv"
	"os"
	"encoding/csv"
	"time"
	"github.com/sgjp/LindaExperimentServerMC/multichain"
)
var taskDurationFile = "/Users/jsanchez/workspace/src/github.com/sgjp/LindaExperimentServerMC/TaskDuration.csv"

var primeNumsQty int

var startTime time.Time

var flag = true

var resultQty int

func StartServer() {

	primeNumsQty = 5000
	resultQty = 0

	log.Fatal(coap.ListenAndServeMulticast("udp", "224.0.1.187:5683",
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			//log.Printf("Got message path=%q: PayLoad: %#v from %v Code: %v", m.Path(), string(m.Payload), a, m.Code)
			if len(m.Path()) > 0 {

				switch m.Path()[0] {

				case "in":
					//log.Printf("COAP message received!")
					res := inTuple(m)
					//log.Printf("COAP message answered: %v!",string(res.Payload))
					return res
				case "out":
					//log.Printf("COAP message received!")
					res := outTuple(m)
					//log.Printf("COAP message answered: %v!",string(res.Payload))
					return res
				default:
					res := notFoundHandler(m)
					return res

				}
			} else {
				res := notFoundHandler(m)
				return res
			}
			return nil
		}),"en1"))


}
func inTuple(m *coap.Message) *coap.Message {

	log.Printf("Searching Tuple: %v",(string(m.Payload)))
	t1 := tupleSpace.Take(string(m.Payload))
	log.Printf("Tuple found: %v ",t1)

	payload := itemToPayload(t1)

	res := &coap.Message{
		Type:      coap.Acknowledgement,
		Code:      coap.Content,
		MessageID: m.MessageID,
		Token:     m.Token,
		Payload:   payload,
	}
	res.SetOption(coap.ContentFormat, coap.TextPlain)
	return res
}


func outTuple(m *coap.Message) *coap.Message {
	item := payloadToItem(m.Payload)


	//Start counting time when the first W tuple comes
	if flag && item.Key=="W"{
		startTime = time.Now()
		log.Println("First W tuple arrived!")
		flag = false
	}
	go tupleSpace.Write(item)
	log.Printf("Outing tuple: %v.",item)

	if item.Key=="R"{
		resultQty++
		log.Printf("ResultQTY %v.",resultQty)
		if resultQty==primeNumsQty{
			elapsed := time.Since(startTime)
			saveTaskDuration(int64(elapsed/time.Millisecond),primeNumsQty)
			log.Println("Last R tuple arrived!")
		}

	}

	res := &coap.Message{
		Type:      coap.Acknowledgement,
		Code:      coap.Created,
		MessageID: m.MessageID,
		Token:     m.Token,
		Payload:   []byte(string("1")),
	}
	res.SetOption(coap.ContentFormat, coap.TextPlain)
	return res
}


func payloadToItem(payload []byte) multichain.Item{
	var data []string
	payloadString := string(payload)

	data = strings.Split(payloadString,",")

	item := multichain.Item{Key:strings.Replace(data[0],"\"","",1),Data:strings.Replace(data[1],"\"","",1)}

	return item

}

func itemToPayload(item multichain.Item) []byte{
	if (item.Key==""){
		return []byte(item.Data)
	}
	return []byte(item.Key+","+item.Data)

}


func notFoundHandler(m *coap.Message) *coap.Message {

	res := &coap.Message{
		Type:      coap.Acknowledgement,
		Code:      coap.NotFound,
		MessageID: m.MessageID,
		Token:     m.Token,
		Payload:   []byte("4.05"),
	}
	res.SetOption(coap.ContentFormat, coap.TextPlain)
	return res

}


func saveTaskDuration(elapsed int64, qty int){
	record := []string{
		strconv.Itoa(qty), strconv.FormatInt(elapsed,10)}

	file, er := os.OpenFile(taskDurationFile, os.O_RDWR|os.O_APPEND, 0666)

	if er != nil {
		log.Fatal(er)
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	err := writer.Write(record)


	if err != nil {
		log.Fatal(er)
	}

	defer writer.Flush()
}
