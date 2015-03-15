package client

import (
	"github.com/davecgh/go-spew/spew"

	"encoding/json"
)

type Channel struct {
	Client        *Client
	SourceId      string
	DestinationId string
	Namespace     string
	RequestId     int
}

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
	spew.Dump(msg)
	return c.Client.Send(&msg)
}
