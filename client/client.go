package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"

	"github.com/golang/protobuf/proto"
)

//go:generate protoc --go_out=. cast_channel.proto

type Client struct {
	Host   string
	Port   int
	Stream *PacketStream
}

func (c *Client) Connect() error {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port), &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return fmt.Errorf("Failed to connect:%s", err)
	}

	c.Stream = NewPacketStream(conn)

	channel := Channel{
		Client:        c,
		SourceId:      "sender-0",
		DestinationId: "receiver-0",
		Namespace:     "urn:x-cast:com.google.cast.tp.connection",
	}
	err = channel.Send(&Payload{
		Type: "CONNECT",
	})
	if err != nil {
		return fmt.Errorf("Failed to connect: %s", err)
	}

	go func() {
		for {
			packet := c.Stream.Read()
			msg := CastMessage{}
			if err := proto.Unmarshal(packet, &msg); err != nil {
				log.Fatalln("Failed to unmarshal CastMessage:", err)
			}
			spew.Dump(msg)

			var payload Payload
			if err := json.Unmarshal([]byte(*msg.PayloadUtf8), &payload); err != nil {
				log.Fatalln("Failed to unmarshal payload:", err)
			}
			spew.Dump(payload)
		}
	}()

	return nil
}

func (c *Client) Send(msg *CastMessage) error {
	proto.SetDefaults(msg)
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = c.Stream.Write(data)
	return err
}

type Payload struct {
	Type      string `json:"type"`
	RequestId *int   `json:"requestId,omitempty"`
}
