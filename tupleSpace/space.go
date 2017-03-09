package tupleSpace

import (
	//"code.google.com/p/go-uuid/uuid"

	"github.com/sgjp/LindaExperimentServerMC/multichain"
	"log"
)

var UsedTxIds []string

var streamCache multichain.Stream

var streamName = "js158"

func Take(key string) multichain.Item {
	if len(streamCache.Response) <= 0 || len(streamCache.Response) == len(UsedTxIds) {
		log.Printf("Getting streamCache")
		streamCache = multichain.GetStream(streamName)
	}
	if len(streamCache.Response) <= 0 {
		log.Printf("StreamCache len<=0")
		return multichain.Item{Key:"", Data:"0"}
	}
	i:=0
	for _, element := range streamCache.Response {
		i++
		if (element.Key == key) {
			if !contains(UsedTxIds, element.Txid) {
				UsedTxIds = append(UsedTxIds, element.Txid)
				log.Printf("Element found: %v",element)
				return multichain.Item{Key:element.Key, Data:element.Data}
			}
		}

	}
	log.Printf("No element found, run %v times. LEN streamCache: %v, LEN UsedTXIDS: %v",i,len(streamCache.Response), len(UsedTxIds))
	return multichain.Item{Key:"", Data:"0"}
}

func Write(item multichain.Item) {
	multichain.AddItemToStream(item.Key, item.Data, streamName)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}