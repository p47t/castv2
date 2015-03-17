package client

type receiverController struct {
	channel *Channel
	reqId   int
}

func NewReceiverController(client *Client, sourceId, destinationId string) *receiverController {
	return &receiverController{
		channel: client.NewChannel(sourceId, destinationId, "urn:x-cast:com.google.cast.tp.receiver"),
		reqId:   1,
	}
}

func (c *receiverController) GetStatus() error {
	c.channel.Send(&LaunchPayload{
		Payload: Payload{
			Type:      "LAUNCH",
			RequestId: &c.reqId,
		},
		AppId: "YouTube",
	})
	c.reqId++

	c.channel.Send(&Payload{
		Type: "GET_STATUS",
	})
	return nil
}
