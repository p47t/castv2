package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type MediaController struct {
	interval       time.Duration
	channel        *Channel
	DestinationID  string
	MediaSessionID int
}

var getMediaStatus = Payload{Type: "GET_STATUS"}
var commandMediaPlay = Payload{Type: "PLAY"}
var commandMediaPause = Payload{Type: "PAUSE"}
var commandMediaStop = Payload{Type: "STOP"}

type MediaCommand struct {
	Payload
	MediaSessionID int `json:"mediaSessionId"`
}

func NewMediaController(client *Client, sourceId, destinationID string) *MediaController {
	controller := &MediaController{
		channel:       client.NewChannel(sourceId, destinationID, "urn:x-cast:com.google.cast.media"),
		DestinationID: destinationID,
	}

	controller.channel.OnMessage("MEDIA_STATUS", func(message *CastMessage) {
		controller.onStatus(message)
	})

	return controller
}

func (c *MediaController) SetDestinationID(id string) {
	c.channel.destinationId = id
	c.DestinationID = id
}

func (c *MediaController) onStatus(message *CastMessage) ([]*MediaStatus, error) {
	spew.Dump("Got media status message", message)

	response := &MediaStatusResponse{}
	err := json.Unmarshal([]byte(*message.PayloadUtf8), response)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal status message:%s - %s", err, *message.PayloadUtf8)
	}

	return response.Status, nil
}

type MediaStatusResponse struct {
	Payload
	Status []*MediaStatus `json:"status,omitempty"`
}

type MediaStatus struct {
	Payload
	MediaSessionID         int                    `json:"mediaSessionId"`
	PlaybackRate           float64                `json:"playbackRate"`
	PlayerState            string                 `json:"playerState"`
	CurrentTime            float64                `json:"currentTime"`
	SupportedMediaCommands int                    `json:"supportedMediaCommands"`
	Volume                 *Volume                `json:"volume,omitempty"`
	CustomData             map[string]interface{} `json:"customData"`
	IdleReason             string                 `json:"idleReason"`
}

func (c *MediaController) GetStatus(timeout time.Duration) ([]*MediaStatus, error) {
	spew.Dump("getting media Status")

	message, err := c.channel.Request(&getMediaStatus)
	if err != nil {
		return nil, fmt.Errorf("Failed to get receiver status: %s", err)
	}

	spew.Dump("got media Status", message)
	return c.onStatus(message)
}

func (c *MediaController) Play(timeout time.Duration) (*CastMessage, error) {
	message, err := c.channel.Request(&MediaCommand{commandMediaPlay, c.MediaSessionID})
	if err != nil {
		return nil, fmt.Errorf("Failed to send play command: %s", err)
	}

	return message, nil
}

func (c *MediaController) Pause(timeout time.Duration) (*CastMessage, error) {
	message, err := c.channel.Request(&MediaCommand{commandMediaPause, c.MediaSessionID})
	if err != nil {
		return nil, fmt.Errorf("Failed to send pause command: %s", err)
	}

	return message, nil
}

func (c *MediaController) Stop(timeout time.Duration) (*CastMessage, error) {
	message, err := c.channel.Request(&MediaCommand{commandMediaStop, c.MediaSessionID})
	if err != nil {
		return nil, fmt.Errorf("Failed to send stop command: %s", err)
	}

	return message, nil
}
