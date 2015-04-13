package client

//go:generate protoc --go_out=. cast_channel.proto

type Channel struct {
	client        *Client
	destinationId string
	namespace     string
}

// Send converts specified payload to JSON and sends wrapped message
func (c *Channel) Send(payload interface{}) error {
	return c.client.Send(c.destinationId, c.namespace, payload)
}

func (c *Channel) Request(req Request) (Response, error) {
	return c.client.Request(c.destinationId, c.namespace, req)
}
