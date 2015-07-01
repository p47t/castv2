package client

// Payload represents JSON payload in a CastMessage
type Payload struct {
	Type      string `json:"type"`
	RequestId *int   `json:"requestId,omitempty"`
}

func (p *Payload) setRequestId(reqId int) {
	p.RequestId = &reqId
}

type Request interface {
	setRequestId(reqId int)
}
