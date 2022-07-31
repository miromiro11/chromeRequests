package models

type Response struct {
	StatusCode int               `json:"StatusCode,omitempty"`
	Body       string            `json:"Body,omitempty"`
	Cookies    map[string]string `json:"Cookies,omitempty"`
	Headers    map[string]string `json:"Headers,omitempty"`
	Url        string            `json:"Url,omitempty"`
	Error      string            `json:"Error,omitempty"`
	SessionId  string            `json:"SessionId,omitempty"`
}
