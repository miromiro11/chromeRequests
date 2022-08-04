package models

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Body       string            `json:"body,omitempty"`
	Cookies    map[string]string `json:"cookies,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Url        string            `json:"url,omitempty"`
	Error      string            `json:"error,omitempty"`
	SessionId  string            `json:"sessionId,omitempty"`
}
