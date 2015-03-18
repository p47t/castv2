package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/proto"
)

type Client struct {
	host          string
	port          int
	packetsStream *PacketStream
	channels      []*Channel
}

func NewClient(host string, port int) (*Client, error) {
	c := &Client{
		host: host,
		port: port,
	}
	hostAddr := fmt.Sprintf("%s:%d", c.host, c.port)
	log.Println("Dialing to:", hostAddr)
	conn, err := tls.Dial("tcp", hostAddr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect:%s", err)
	}

	c.packetsStream = NewPacketStream(conn)

	go func() {
		for {
			packet := c.packetsStream.Read()
			msg := CastMessage{}
			if err := proto.Unmarshal(packet, &msg); err != nil {
				log.Fatalln("Failed to unmarshal CastMessage:", err)
			}
			log.Println("Got message:", spew.Sdump(msg))

			var payload Payload
			if err := json.Unmarshal([]byte(*msg.PayloadUtf8), &payload); err != nil {
				log.Fatalln("Failed to unmarshal payload:", err)
			}

			for _, ch := range c.channels {
				ch.OnMessage(&msg, &payload)
			}
		}
	}()

	return c, nil
}

func (c *Client) NewChannel(sourceId, destinationId, namespace string) *Channel {
	ch := Channel{
		client:        c,
		sourceId:      sourceId,
		destinationId: destinationId,
		namespace:     namespace,
		listeners:     make(map[string]func(*CastMessage)),
	}
	c.channels = append(c.channels, &ch)
	return &ch
}

func (c *Client) Send(msg *CastMessage) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = c.packetsStream.Write(data)
	log.Println("Send message:", spew.Sdump(msg))

	return err
}
