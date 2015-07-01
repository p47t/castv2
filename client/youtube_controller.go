package client

import (
	"encoding/json"
	"fmt"
	"log"
)

type youtubeController struct {
	channel *Channel
}

func NewYouTubeController(client *Client, sourceId, destinationId string) *youtubeController {
	return &youtubeController{
		channel: client.NewChannel(sourceId, destinationId, "urn:x-cast:com.google.youtube.mdx"),
	}
}

func MustMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatalln("Failed to marshal:", err)
	}
	return string(b)
}

func (c *youtubeController) Connect() error {
	err := c.channel.Send(&Payload{
		Type: "CONNECT",
	})
	if err != nil {
		return fmt.Errorf("Failed to connect: %s", err)
	}
	return nil
}

func (c *youtubeController) Load(videoID string) error {
	c.channel.Send(&flingVideoPayload{
		Payload: Payload{
			Type: "flingVideo",
		},
		Data: flingVideoData{
			CurrentTime: 0,
			VideoID:     videoID,
		},
	})
	return nil
}

func (c *youtubeController) GetStatus() error {
	_, err := c.channel.Request(&Payload{
		Type: "GET_STATUS",
	})

	return err
}

type (
	flingVideoPayload struct {
		Payload
		Data flingVideoData `json:"data"`
	}

	flingVideoData struct {
		CurrentTime int    `json:"currentTime"`
		VideoID     string `json:"videoId"`
	}
)
