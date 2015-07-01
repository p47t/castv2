package client

type connectionController struct {
	channel *Channel
}

var connect = &Payload{Type: "CONNECT"}
var close = &Payload{Type: "CLOSE"}

func NewConnectionController(client *Client, sourceId, destinationId string) *connectionController {
	return &connectionController{
		channel: client.NewChannel(sourceId, destinationId, "urn:x-cast:com.google.cast.tp.connection"),
	}
}

func (c *connectionController) Connect() error {
	return c.channel.Send(connect)
}

func (c *connectionController) Close() error {
	return c.channel.Send(close)
}
