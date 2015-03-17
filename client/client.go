package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
)

type Client struct {
	Host   string
	Port   int
	Stream *PacketStream
}

func NewClient(host string, port int) (*Client, error) {
	c := &Client{
		Host: host,
		Port: port,
	}
	hostAddr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	log.Println("Dialing to:", hostAddr)
	conn, err := tls.Dial("tcp", hostAddr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect:%s", err)
	}

	c.Stream = NewPacketStream(conn)

	go func() {
		for {
			packet := c.Stream.Read()
			msg := CastMessage{}
			if err := proto.Unmarshal(packet, &msg); err != nil {
				log.Fatalln("Failed to unmarshal CastMessage:", err)
			}
			log.Println("Got message:", msg)

			var payload Payload
			if err := json.Unmarshal([]byte(*msg.PayloadUtf8), &payload); err != nil {
				log.Fatalln("Failed to unmarshal payload:", err)
			}
		}
	}()

	return c, nil
}

func (c *Client) NewChannel(sourceId, destinationId, namespace string) *Channel {
	ch := Channel{
		Client:        c,
		SourceId:      sourceId,
		DestinationId: destinationId,
		Namespace:     namespace,
	}
	return &ch
}

func (c *Client) Send(msg *CastMessage) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = c.Stream.Write(data)
	log.Println("Send message:", msg)
	return err
}
