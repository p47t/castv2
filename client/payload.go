package client

// Payload represents JSON payload in a CastMessage
type Payload struct {
	Type      string `json:"type"`
	RequestId *int   `json:"requestId,omitempty"`
}

func (p *Payload) setRequestId(reqId int) {
	p.RequestId = &reqId
}

func (p *Payload) getRequestId() int {
	return *p.RequestId
}

type LaunchPayload struct {
	Payload
	AppId string `json:"appId"`
}

type requestIdCarrier interface {
	setRequestId(reqId int)
	getRequestId() int
}
