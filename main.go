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

	cc := client.NewConnectionController(c, "sender-0", "receiver-0")
	err = cc.Connect()

	client.NewHeartBeatController(c, "sender-0", "receiver-0")

	rc := client.NewReceiverController(c, "sender-0", "receiver-0")
	rc.GetStatus()
	rc.Launch("YouTube")

	<-chan interface{}(nil)
}
