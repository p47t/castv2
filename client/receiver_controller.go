package client

import (
	"encoding/json"
	"log"
)

type receiverController struct {
	channel *Channel
}

var getStatus = &Payload{Type: "GET_STATUS"}

func NewReceiverController(client *Client, sourceId, destinationId string) *receiverController {
	c := &receiverController{
		channel: client.NewChannel(sourceId, destinationId, "urn:x-cast:com.google.cast.receiver"),
	}

	c.channel.OnMessage("RECEIVER_STATUS", c.onStatus)

	return c
}

type StatusResponse struct {
	Payload
	Status *ReceiverStatus `json:"status,omitempty"`
}

type ReceiverStatus struct {
	Applications  []*ApplicationSession `json:"applications"`
	IsActiveInput *bool                 `json:"isActiveInput,omitempty"`
	IsStandBy     *bool                 `json:"isStandBy,omitempty"`
	Volume        *Volume               `json:"volume,omitempty"`
}

type ApplicationSession struct {
	AppID       *string      `json:"appId,omitempty"`
	DisplayName *string      `json:"displayName,omitempty"`
	Namespaces  []*Namespace `json:"namespaces"`
	SessionID   *string      `json:"sessionId,omitempty"`
	StatusText  *string      `json:"statusText,omitempty"`
	TransportId *string      `json:"transportId,omitempty"`
}

type Namespace struct {
	Name string `json:"name"`
}

type Volume struct {
	Level *float64 `json:"level,omitempty"`
	Muted *bool    `json:"muted,omitempty"`
}

func (c *receiverController) onStatus(msg *CastMessage) {
	status := StatusResponse{}
	err := json.Unmarshal([]byte(*msg.PayloadUtf8), &status)
	if err != nil {
		log.Printf("Failed to unmarshal status message:%s - %s", err, *msg.PayloadUtf8)
		return
	}
}

type LaunchPayload struct {
	Payload
	AppId string `json:"appId"`
}

func (c *receiverController) GetStatus() (*ReceiverStatus, error) {
	msg, err := c.channel.Request(getStatus)

	status := StatusResponse{}
	err = json.Unmarshal([]byte(*msg.PayloadUtf8), &status)
	if err != nil {
		log.Printf("Failed to unmarshal status message:%s - %s", err, *msg.PayloadUtf8)
		return nil, err
	}

	return status.Status, nil
}

func (c *receiverController) Launch(appId string) error {
	_, err := c.channel.Request(&LaunchPayload{
		Payload: Payload{
			Type: "LAUNCH",
		},
		AppId: appId,
	})
	return err
}
