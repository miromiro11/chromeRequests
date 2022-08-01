package models

type SessionParameters struct {
	SessionId   string             `json:"session"`
	RequestType string             `json:"requestType"`
	Parameters  *RequestParameters `json:"parameters"`
	Proxy       string             `json:"proxy"`
}
