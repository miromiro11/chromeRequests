package utils

import "C"
import (
	"chromeRequests/models"
	"encoding/json"
	http "github.com/saucesteals/fhttp"
	"net/url"
)

func CreateCResponse(resp *models.Response) *C.char {
	errorJson, _ := json.Marshal(resp)
	return C.CString(string(errorJson))
}

func CreateTransport(proxy string) (*http.Transport, error) {
	if len(proxy) != 0 {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		return &http.Transport{Proxy: http.ProxyURL(proxyUrl)}, nil
	} else {
		return &http.Transport{}, nil
	}
}
