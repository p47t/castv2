package client

import (
	"time"
)

type heartBeatController struct {
	ticker  *time.Ticker
	channel *Channel
}

var ping = &Payload{Type: "PING"}
var pong = &Payload{Type: "PONG"}

func NewHeartBeatController(client *Client, sourceId, destinationId string) *heartBeatController {
	c := &heartBeatController{
		channel: client.NewChannel(sourceId, destinationId, "urn:x-cast:com.google.cast.tp.heartbeat"),
	}

	// Returns a PONG for incoming PING
	c.channel.OnMessage("PING", func(_ *CastMessage) {
		c.channel.Send(pong)
	})

	return c
}

func (c *heartBeatController) Start() {
	if c.ticker != nil {
		c.Stop()
	}

	// Ping destination every 5 minutes
	c.ticker = time.NewTicker(5 * time.Second)
	go func() {
		for {
			c.channel.Send(ping)
			<-c.ticker.C
		}
	}()
}

func (c *heartBeatController) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
		c.ticker = nil
	}
}
