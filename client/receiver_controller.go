package client

type receiverController struct {
	channel *Channel
}

func NewReceiverController(client *Client, destinationId string) *receiverController {
	return &receiverController{
		channel: client.NewChannel(destinationId, "urn:x-cast:com.google.cast.receiver"),
	}
}

func (c *receiverController) GetStatus() error {
	c.channel.Request(&Payload{
		Type: "GET_STATUS",
	})

	return nil
}

func (c *receiverController) Launch(appId string) error {
	c.channel.Request(&LaunchPayload{
		Payload: Payload{
			Type: "LAUNCH",
		},
		AppId: appId,
	})
	return nil
}
