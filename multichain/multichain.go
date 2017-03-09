package multichain

import (
	"net/url"
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"encoding/hex"
	"bytes"
)

var mcUrl = "http://10.81.4.68:8080/streams"



type Stream struct{
	Response []StreamItem `json:"response"`
}

type StreamItem struct{
	Publishers []string `json:"publishers"`
	Confirmations int `json:"confirmations"`
	Data string `json:"data"`
	Key string `json:"key`
	Txid string `json:"txid"`
}

type Item struct{
	Key string `json:"key"`
	Data string `json:"data"`
}

func GetStream(streamName string) Stream{

	safeStreamName := url.QueryEscape(streamName)

	url := fmt.Sprintf(mcUrl+"/%s",safeStreamName)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		//return
	}


	client := &http.Client{}


	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		//return
	}

	defer resp.Body.Close()
	var record Stream

	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	for index,element := range record.Response{
		value,_ :=hex.DecodeString(element.Data)
		element.Data=string(value)
		record.Response[index]=element
	}

	return record

}

func AddItemToStream(key, data, streamName string){
	item := Item{Key:key,Data:data}
	itemByte,_:=json.Marshal(item)

	safeStreamName := url.QueryEscape(streamName)

	url := fmt.Sprintf(mcUrl+"/%s",safeStreamName)

	// Build the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(itemByte))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}


	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()


}

