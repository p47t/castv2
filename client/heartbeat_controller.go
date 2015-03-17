package client

import (
	"time"
)

type heartBeatController struct {
	ticker  *time.Ticker
	channel *Channel
}

func NewHeartBeatController(client *Client, sourceId, destinationId string) *heartBeatController {
	return &heartBeatController{
		channel: client.NewChannel(sourceId, destinationId, "urn:x-cast:com.google.cast.tp.heartbeat"),
	}
}

func (c *heartBeatController) Start() error {
	c.ticker = time.NewTicker(5 * time.Second)
	go func() {
		for {
			c.channel.Send(&Payload{
				Type: "PING",
			})
			<-c.ticker.C
		}
	}()
	return nil
}

func (c *heartBeatController) OnPing() {
	c.channel.Send(&Payload{
		Type: "PONG",
	})
}
