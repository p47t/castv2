package client

import "encoding/json"

//go:generate protoc --go_out=. cast_channel.proto

type Channel struct {
	Client        *Client
	SourceId      string
	DestinationId string
	Namespace     string
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
		SourceId:        &c.SourceId,
		DestinationId:   &c.DestinationId,
		Namespace:       &c.Namespace,
		PayloadType:     CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadStr,
	}
	return c.Client.Send(&msg)
}
