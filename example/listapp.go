package main

import (
	"encoding/json"
	"hlcient/hlclient"
	"log"
)

func main() {
	hl := hlclient.NewClient("127.0.0.1:32876")

	err := hl.Connect()
	if err != nil {
		log.Fatalf("connect err: %v", err)
	}

	msgData, err := json.Marshal(map[string]interface{}{
		"index": 1,
	})
	if err != nil {
		log.Fatalf("marshal err: %v", err)
	}

	resp, err := hl.SendMessage(']', msgData, true)
	if err != nil {
		log.Fatalf("send err: %v", err)
	} else {
		log.Printf("resp: %v", resp)
		log.Printf("resp.code: %d", resp.Code)
		log.Printf("resp.data: %s", string(resp.Data))
	}
}
