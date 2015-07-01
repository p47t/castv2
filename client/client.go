package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

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
		log.Println("Failed to dial:", err)
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}

	c.packetsStream = NewPacketStream(conn)
	go c.dispatchResponses()

	// Add a new heart beat controller automatically
	hbc := NewHeartBeatController(c, "sender-0", "receiver-0")
	hbc.Start()

	return c, nil
}

func (c *Client) dispatchResponses() {
	for {
		packet := c.packetsStream.Read()
		msg := CastMessage{}
		if err := proto.Unmarshal(packet, &msg); err != nil {
			log.Fatalln("Failed to unmarshal CastMessage:", err)
		}
		log.Printf("Recv: S=%s, D=%s, NS=%s, %s", *msg.SourceId, *msg.DestinationId, *msg.Namespace, *msg.PayloadUtf8)

		var headers Payload
		if err := json.Unmarshal([]byte(*msg.PayloadUtf8), &headers); err != nil {
			log.Fatalln("Failed to unmarshal headers:", err)
		}

		for _, channel := range c.channels {
			channel.message(&msg, &headers)
		}
	}
}

func (c *Client) NewChannel(sourceId, destinationId, namespace string) *Channel {
	ch := Channel{
		client:        c,
		sourceId:      sourceId,
		destinationId: destinationId,
		namespace:     namespace,
		inFlight:      make(map[int]chan Response),
		listeners:     make([]channelListener, 0),
	}
	c.channels = append(c.channels, &ch)
	return &ch
}

func (c *Client) sendCastMessage(msg *CastMessage) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = c.packetsStream.Write(data)

	return err
}

// Send converts specified payload to JSON and sends wrapped message
func (c *Client) Send(sourceId, destinationId, namespace string, payload interface{}) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	payloadStr := string(payloadJson)
	msg := CastMessage{
		ProtocolVersion: CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &sourceId,
		DestinationId:   &destinationId,
		Namespace:       &namespace,
		PayloadType:     CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadStr,
	}
	log.Println("Send:", payloadStr)
	return c.sendCastMessage(&msg)
}
