package client

import "fmt"

type connectionController struct {
	channel *Channel
}

func NewConnectionController(client *Client, destinationId string) *connectionController {
	return &connectionController{
		channel: client.NewChannel(destinationId, "urn:x-cast:com.google.cast.tp.connection"),
	}
}

func (c *connectionController) Connect() error {
	err := c.channel.Send(&Payload{
		Type: "CONNECT",
	})
	if err != nil {
		return fmt.Errorf("Failed to connect: %s", err)
	}
	return nil
}
