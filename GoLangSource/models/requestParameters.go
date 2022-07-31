package models

type RequestParameters struct {
	URL       string            `json:"url"`
	Proxy     string            `json:"proxy"`
	Headers   map[string]string `json:"headers"`
	Form      map[string]string `json:"FORM"`
	Json      string            `json:"JSON"`
	Cookies   map[string]string `json:"cookies"`
	Redirects bool              `json:"redirects"`
}
