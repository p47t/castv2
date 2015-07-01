package main

import (
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"

	"simplypatrick.com/castv2/client"
)

func main() {
	info := <-client.SearchChromecast()
	log.Printf("Found Chromecast at %s:%d\n", info.Host, info.Port)

	c, err := client.NewClient(info.Host, info.Port)
	if err != nil {
		log.Fatalln("Failed to connect:", err)
	}

	cc := client.NewConnectionController(c, "sender-0", "receiver-0")
	err = cc.Connect()

	rc := client.NewReceiverController(c, "sender-0", "receiver-0")
	rc.Launch("CC1AD845")

	time.Sleep(5 * time.Second)

	rstatus, _ := rc.GetStatus()
	spew.Dump(rstatus)

	mc := client.NewMediaController(c, "sender-0", *rstatus.Applications[0].TransportId)
	mstatus, _ := mc.GetStatus(0)
	spew.Dump(mstatus)

	// yc := client.NewYouTubeController(c, "sender-0", *rstatus.Applications[0].TransportId)
	// yc.Connect()
	// yc.Load("9bZkp7q19f0")
	// yc.GetStatus()

	<-chan interface{}(nil)
}
