package client

type Response interface{}

type Status struct {
	requestId    int           `json:"requestId,omitempty"`
	applications []Application `json:"applications"`
}

type Application struct {
	appId       string   `json:"appId"`
	displayName string   `json:"displayName"`
	sessionId   string   `json:"sessionId"`
	statusText  string   `json:"statusText"`
	transportId string   `json:"transportId"`
	namespaces  []string `json:"namespaces"`
}
