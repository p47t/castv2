package main

import (
	"log"

	"simplypatrick.com/castv2/client"
)

func main() {
	launchYoutuneOnChromecast("9bZkp7q19f0")
}

func launchYoutuneOnChromecast(videoID string) {
	info := <-client.SearchChromecast()
	log.Printf("Found Chromecast at %s:%d\n", info.Host, info.Port)

	c, err := client.NewClient(info.Host, info.Port)
	if err != nil {
		log.Fatalln("Failed to connect:", err)
	}

	cc := client.NewConnectionController(c, "receiver-0")
	err = cc.Connect()

	rc := client.NewReceiverController(c, "receiver-0")
	rc.GetStatus()

	yc := client.NewYouTubeController(c, "receiver-0")
	yc.Connect()
	yc.Load(videoID)
	yc.GetStatus()

	<-chan interface{}(nil)
}
