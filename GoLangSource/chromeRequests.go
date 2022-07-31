package chromeRequests

import (
	"C"
	"bytes"
	"chromeRequests/models"
	"chromeRequests/utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	http "github.com/saucesteals/fhttp"
	"github.com/saucesteals/mimic"
	"io/ioutil"
	"net/url"
)

var Sessions = make(map[string]models.Session)
var userAgent = fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", m.Version())
var latestVersion = mimic.MustGetLatestVersion(mimic.PlatformWindows)
var m, _ = mimic.Chromium(mimic.BrandChrome, latestVersion)
var cleanTransport = &http.Transport{}

//export changeProxy
func changeProxy(params *C.char) *C.char {
	var sessionParameters models.SessionParameters
	err := json.Unmarshal([]byte(C.GoString(params)), &sessionParameters)
	if err != nil {
		return utils.CreateCResponse(&models.Response{Error: err.Error()})
	}
	proxy := sessionParameters.Proxy
	sessionId := sessionParameters.Session

	if session, exists := Sessions[sessionId]; exists {
		transport, err := utils.CreateTransport(proxy)
		if err != nil {
			return utils.CreateCResponse(&models.Response{Error: err.Error()})
		}

		session.Client.Transport = m.ConfigureTransport(transport)
	}

	return utils.CreateCResponse(&models.Response{})
}

//export createSession
func createSession(cProxy *C.char) *C.char {
	proxy := C.GoString(cProxy)
	sessionId := uuid.NewString()

	transport, err := utils.CreateTransport(proxy)
	if err != nil {
		return utils.CreateCResponse(&models.Response{Error: err.Error()})
	}

	Sessions[sessionId] = models.Session{
		Client:    &http.Client{Transport: m.ConfigureTransport(transport)},
		Headers:   make(map[string]string),
		Randomize: false,
		Cookies:   make(map[string]string),
	}
	return utils.CreateCResponse(&models.Response{SessionId: sessionId})
}

//export closeSession
func closeSession(uuid *C.char) *C.char {
	if session, exists := Sessions[C.GoString(uuid)]; exists {
		session.Client.CloseIdleConnections()
	} else {
		return utils.CreateCResponse(&models.Response{Error: "session does not exists"})
	}

	return utils.CreateCResponse(&models.Response{})
}

//export request
func request(cParams *C.char) *C.char {
	var client *http.Client
	var req *http.Request
	params := C.GoString(cParams)
	data := models.SessionParameters{}
	err := json.Unmarshal([]byte(params), &data)
	if err != nil {
		return utils.CreateCResponse(&models.Response{Error: err.Error()})
	}

	if data.Session != "" {
		client = Sessions[data.Session].Client
	} else {
		transport, err := utils.CreateTransport("")
		if err != nil {
			return utils.CreateCResponse(&models.Response{Error: err.Error()})
		}

		client = &http.Client{
			Transport: m.ConfigureTransport(transport),
		}
	}
	if data.RequestType == "GET" {
		req, err = http.NewRequest("GET", data.Parameters.URL, nil)
		if err != nil {
			return utils.CreateCResponse(&models.Response{Error: err.Error()})
		}
	} else if data.RequestType == "POST" || data.RequestType == "PUT" {
		req, err = http.NewRequest(data.RequestType, data.Parameters.URL, nil)
		if err != nil {
			return utils.CreateCResponse(&models.Response{Error: err.Error()})
		}

		if len(data.Parameters.Form) != 0 {
			formData := url.Values{}
			for key, value := range data.Parameters.Form {
				formData.Add(key, value)
			}
			req, err = http.NewRequest(data.RequestType, data.Parameters.URL, bytes.NewBufferString(formData.Encode()))
			if err != nil {
				return utils.CreateCResponse(&models.Response{Error: err.Error()})
			}
		} else if data.Parameters.Json != "" {
			req, err = http.NewRequest(data.RequestType, data.Parameters.URL, bytes.NewBuffer([]byte(data.Parameters.Json)))
			if err != nil {
				return utils.CreateCResponse(&models.Response{Error: err.Error()})
			}
		}
	}
	req.Header = http.Header{
		"sec-ch-ua":          {m.ClientHintUA()},
		"rtt":                {"50"},
		"sec-ch-ua-mobile":   {"?0"},
		"user-agent":         {userAgent},
		"accept":             {"text/html,*/*"},
		"x-requested-with":   {"XMLHttpRequest"},
		"downlink":           {"3.9"},
		"ect":                {"4g"},
		"sec-ch-ua-platform": {`"Windows"`},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"en,en_US;q=0.9"},
		http.HeaderOrderKey: {
			"sec-ch-ua", "rtt", "sec-ch-ua-mobile",
			"user-agent", "accept", "x-requested-with",
			"downlink", "ect", "sec-ch-ua-platform",
			"sec-fetch-site", "sec-fetch-mode", "sec-fetch-dest",
			"accept-encoding", "accept-language",
		},
		http.PHeaderOrderKey: m.PseudoHeaderOrder(),
	}
	if data.Parameters.Proxy != "" {
		transport, err := utils.CreateTransport(data.Parameters.Proxy)
		if err != nil {
			return utils.CreateCResponse(&models.Response{Error: err.Error()})
		}
		client.Transport = m.ConfigureTransport(transport)
	}
	if data.Parameters.Headers != nil {
		for k, v := range data.Parameters.Headers {
			req.Header.Set(k, v)
		}
	}

	if data.Parameters.Json != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	if len(data.Parameters.Form) != 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	for k, v := range data.Parameters.Cookies {
		req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
	if !data.Parameters.Redirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return utils.CreateCResponse(&models.Response{Error: err.Error()})
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return utils.CreateCResponse(&models.Response{Error: err.Error()})
	}

	headersMap := make(map[string]string)
	for key, value := range resp.Header {
		headersMap[key] = value[0]
	}

	cookieMap := make(map[string]string)
	for _, cookie := range resp.Cookies() {
		cookieMap[cookie.Name] = cookie.Value
	}

	client.Transport = cleanTransport
	client.CheckRedirect = nil

	return utils.CreateCResponse(&models.Response{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Cookies:    cookieMap,
		Headers:    headersMap,
		Url:        resp.Request.URL.String(),
	})
}
