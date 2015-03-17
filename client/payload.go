package client

type Payload struct {
	Type      string `json:"type"`
	RequestId *int   `json:"requestId,omitempty"`
}

type LaunchPayload struct {
	Payload
	AppId string `json:"appId"`
}
