package main

import (
	"log"

	"simplypatrick.com/castv2/client"
)

func main() {
	info := <-client.SearchChromecast()

	c, err := client.NewClient(info.Host, info.Port)
	if err != nil {
		log.Fatalln("Failed to connect:", err)
	}

	hc := client.NewHeartBeatController(c, "sender-0", "receiver-0")
	err = hc.Start()

	cc := client.NewConnectionController(c, "sender-0", "receiver-0")
	err = cc.Connect()

	rc := client.NewReceiverController(c, "sender-0", "receiver-0")
	rc.GetStatus()

	<-chan interface{}(nil)
}
