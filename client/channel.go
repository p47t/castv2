package client

import "encoding/json"

//go:generate protoc --go_out=. cast_channel.proto

type Channel struct {
	client        *Client
	sourceId      string
	destinationId string
	namespace     string
	reqId         int
	listeners     map[string]func(*CastMessage)
}

// Request sends request with request ID and wait for response
func (c *Channel) Request(req requestIdCarrier) error {
	c.reqId++
	req.setRequestId(c.reqId)
	if err := c.Send(req); err != nil {
		return err
	}

	// TODO: wait for response
	return nil
}

// Send converts specified payload to JSON and sends wrapped message
func (c *Channel) Send(payload interface{}) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	payloadStr := string(payloadJson)
	msg := CastMessage{
		ProtocolVersion: CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &c.sourceId,
		DestinationId:   &c.destinationId,
		Namespace:       &c.namespace,
		PayloadType:     CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadStr,
	}
	return c.client.Send(&msg)
}

func (c *Channel) OnMessage(msg *CastMessage, payload *Payload) {
	// if msg.GetDestinationId() != "*" && (msg.GetSourceId() != c.destinationId || msg.GetDestinationId() != c.sourceId || msg.GetNamespace() != c.namespace) {
	// 	return
	// }
	if listener, ok := c.listeners[payload.Type]; ok {
		listener(msg)
	}
}

func (c *Channel) Listen(responseType string, callback func(*CastMessage)) {
	c.listeners[responseType] = callback
}
