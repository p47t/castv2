package client

import (
	"fmt"
	"log"
	"sync"
	"time"
)

//go:generate protoc --go_out=. cast_channel.proto

type Channel struct {
	sync.Mutex

	client        *Client
	sourceId      string
	destinationId string
	namespace     string
	inFlight      map[int]chan Response
	listeners     []channelListener
	nextReqId     int
}

type channelListener struct {
	responseType string
	callback     func(*CastMessage)
}

func (c *Channel) NextReqId() int {
	reqId := c.nextReqId
	c.nextReqId++
	return reqId
}

// Send converts specified payload to JSON and sends wrapped message
func (c *Channel) Send(payload interface{}) error {
	return c.client.Send(c.sourceId, c.destinationId, c.namespace, payload)
}

// Request sends request and waits for response
func (c *Channel) Request(req Request) (*CastMessage, error) {
	c.Lock()

	reqId := c.NextReqId()
	req.setRequestId(reqId)

	// Map request ID to result
	result := make(chan Response, 1)
	c.inFlight[reqId] = result

	if err := c.client.Send(c.sourceId, c.destinationId, c.namespace, req); err != nil {
		delete(c.inFlight, reqId)

		c.Unlock()
		return nil, err
	}

	c.Unlock()

	select {
	case response := <-result:
		return response.(*CastMessage), nil
	case <-time.After(30 * time.Second):
		c.Lock()
		delete(c.inFlight, reqId)
		c.Unlock()
		return nil, fmt.Errorf("Timeout:", reqId)
	}
}

func (c *Channel) message(msg *CastMessage, headers *Payload) {
	if *msg.DestinationId != "*" && (*msg.SourceId != c.destinationId || *msg.DestinationId != c.sourceId || *msg.Namespace != c.namespace) {
		return
	}

	c.Lock()
	defer c.Unlock()

	if *msg.DestinationId != "*" && headers.RequestId != nil {
		requester, ok := c.inFlight[*headers.RequestId]
		if !ok {
			log.Println("Unknown reqId", *headers.RequestId)
			return
		}

		// Return msg to requester
		requester <- msg
		delete(c.inFlight, *headers.RequestId)
		return
	}

	if headers.Type == "" {
		log.Println("No message type:", msg)
		return
	}

	for _, l := range c.listeners {
		if l.responseType == headers.Type {
			l.callback(msg)
		}
	}
}

func (c *Channel) OnMessage(responseType string, callback func(*CastMessage)) {
	c.Lock()
	defer c.Unlock()

	c.listeners = append(c.listeners, channelListener{responseType, callback})
}
