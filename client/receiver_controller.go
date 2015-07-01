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
	_, err := c.channel.Request(&Payload{
		Type: "GET_STATUS",
	})

	return err
}

func (c *receiverController) Launch(appId string) error {
	_, err := c.channel.Request(&LaunchPayload{
		Payload: Payload{
			Type: "LAUNCH",
		},
		AppId: appId,
	})
	return err
}
