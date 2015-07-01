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
	name          string
	packetsStream *PacketStream
	channels      []*Channel
	requests      map[int]chan string
	nextReqId     int
}

func NewClient(host string, port int) (*Client, error) {
	c := &Client{
		host:     host,
		port:     port,
		requests: make(map[int]chan string),
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
	return c, nil
}

func (c *Client) dispatchResponses() {
	for {
		packet := c.packetsStream.Read()
		msg := CastMessage{}
		if err := proto.Unmarshal(packet, &msg); err != nil {
			log.Fatalln("Failed to unmarshal CastMessage:", err)
		}

		log.Println("Recv:", *msg.PayloadUtf8)

		// Unmarshal to check requestId
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(*msg.PayloadUtf8), &payload); err != nil {
			log.Fatalln("Failed to unmarshal payload:", err)
		}
		switch requestId := payload["requestId"].(type) {
		case float64:
			reqId := int(requestId)
			fmt.Println("requestId:", reqId)
			if res, ok := c.requests[reqId]; ok {
				// Pass the result to request
				res <- *msg.PayloadUtf8
				delete(c.requests, reqId)
			}
		default:
			switch t := payload["type"].(type) {
			case string:
				if t == "PING" {
					c.Send(*msg.DestinationId, *msg.Namespace, &Payload{Type: "PONG"})
				}
			}
		}
	}
}

func (c *Client) NewChannel(destinationId, namespace string) *Channel {
	ch := Channel{
		client:        c,
		destinationId: destinationId,
		namespace:     namespace,
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
func (c *Client) Send(destinationId, namespace string, payload interface{}) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	payloadStr := string(payloadJson)
	msg := CastMessage{
		ProtocolVersion: CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &c.name,
		DestinationId:   &destinationId,
		Namespace:       &namespace,
		PayloadType:     CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadStr,
	}
	log.Println("Send:", payloadStr)
	return c.sendCastMessage(&msg)
}

func (c *Client) NextReqId() int {
	reqId := c.nextReqId
	c.nextReqId++
	return reqId
}

// Request sends request with request ID and wait for response
func (c *Client) Request(destinationId, namespace string, req Request) (Response, error) {
	reqId := c.NextReqId()
	req.setRequestId(reqId)

	// Map request ID to result
	result := make(chan string, 1)
	c.requests[reqId] = result

	if err := c.Send(destinationId, namespace, req); err != nil {
		return nil, err
	}

	// Wait for result
	return <-result, nil
}

// GetStatus gets receiver status
func (c *Client) GetStatus() (string, error) {
	response, err := c.Request("receiver-0", "urn:x-cast:com.google.cast.receiver", &Payload{})
	return response.(string), err
}
