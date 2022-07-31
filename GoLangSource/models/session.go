package models

import http "github.com/saucesteals/fhttp"

type Session struct {
	Client    *http.Client
	Headers   map[string]string
	Cookies   map[string]string
	Randomize bool
}
