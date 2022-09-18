package main

import (
	"C"
	"bytes"
	"chromeRequests/models"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/google/uuid"

	http "github.com/saucesteals/fhttp"
	"github.com/saucesteals/mimic"
)

var Sessions = make(map[string]models.Session)
var latestVersion = mimic.MustGetLatestVersion(mimic.PlatformWindows)
var m, _ = mimic.Chromium(mimic.BrandChrome, latestVersion)

//export changeProxy
func changeProxy(params *C.char) *C.char {
	var sessionParameters models.SessionParameters
	err := json.Unmarshal([]byte(C.GoString(params)), &sessionParameters)
	if err != nil {
		return createCResponse(&models.Response{Error: err.Error()})
	}
	proxy := sessionParameters.Proxy
	sessionId := sessionParameters.SessionId

	if session, exists := Sessions[sessionId]; exists {
		transport, err := createTransport(proxy)
		if err != nil {
			return createCResponse(&models.Response{Error: err.Error()})
		}

		session.Client.Transport = m.ConfigureTransport(transport)
	}

	return createCResponse(&models.Response{})
}

//export createSession
func createSession(cProxy *C.char) *C.char {
	proxy := C.GoString(cProxy)
	sessionId := uuid.NewString()

	transport, err := createTransport(proxy)
	if err != nil {
		return createCResponse(&models.Response{Error: err.Error()})
	}

	Sessions[sessionId] = models.Session{
		Client:    &http.Client{Transport: m.ConfigureTransport(transport)},
		Headers:   make(map[string]string),
		Randomize: false,
		Cookies:   make(map[string]string),
	}
	return createCResponse(&models.Response{SessionId: sessionId})
}

//export closeSession
func closeSession(uuid *C.char) *C.char {
	if session, exists := Sessions[C.GoString(uuid)]; exists {
		session.Client.CloseIdleConnections()
	} else {
		return createCResponse(&models.Response{Error: "session does not exists"})
	}

	return createCResponse(&models.Response{})
}

//export request
func request(cParams *C.char) *C.char {
	var client *http.Client
	var req *http.Request
	params := C.GoString(cParams)
	data := models.SessionParameters{}
	err := json.Unmarshal([]byte(params), &data)
	if err != nil {
		return createCResponse(&models.Response{Error: err.Error()})
	}

	if data.SessionId != "" {
		client = Sessions[data.SessionId].Client
	} else {
		newClient, err := createClient(data.Parameters.Proxy)
		if err != nil {
			return createCResponse(&models.Response{Error: err.Error()})
		}

		client = newClient
	}

	if data.RequestType == "GET" {
		req, err = http.NewRequest("GET", data.Parameters.URL, nil)
		if err != nil {
			return createCResponse(&models.Response{Error: err.Error()})
		}
	} else if data.RequestType == "POST" || data.RequestType == "PUT" {
		var body io.Reader

		if len(data.Parameters.Form) != 0 {
			formData := url.Values{}
			for key, value := range data.Parameters.Form {
				formData.Add(key, value)
			}

			body = bytes.NewBufferString(formData.Encode())
		} else if data.Parameters.Json != "" {
			body = bytes.NewBuffer([]byte(data.Parameters.Json))
		}

		req, err = http.NewRequest(data.RequestType, data.Parameters.URL, body)
		if err != nil {
			return createCResponse(&models.Response{Error: err.Error()})
		}
	}

	if data.Parameters.Proxy != "" {
		transport, err := createTransport(data.Parameters.Proxy)
		if err != nil {
			return createCResponse(&models.Response{Error: err.Error()})
		}
		client.Transport = m.ConfigureTransport(transport)
	}

	if data.Parameters.Headers != nil {
		req.Header = http.Header{
			http.PHeaderOrderKey: m.PseudoHeaderOrder(),
		}
		for k, v := range data.Parameters.Headers {
			req.Header.Add(k, v)
		}
		req.Header[http.HeaderOrderKey] = data.Parameters.HeaderOrder
		
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
		return createCResponse(&models.Response{Error: err.Error()})
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return createCResponse(&models.Response{Error: err.Error()})
	}

	headersMap := make(map[string]string)
	for key, value := range resp.Header {
		headersMap[key] = value[0]
	}

	cookieMap := make(map[string]string)
	for _, cookie := range resp.Cookies() {
		cookieMap[cookie.Name] = cookie.Value
	}

	client.CheckRedirect = nil
	if data.SessionId == "" {
		client.CloseIdleConnections()
	}

	//fmt.Println( string(body))

	return createCResponse(&models.Response{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Cookies:    cookieMap,
		Headers:    headersMap,
		Url:        resp.Request.URL.String(),
	})
}

func createClient(proxy string) (*http.Client, error) {
	transport, err := createTransport(proxy)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: m.ConfigureTransport(transport),
	}, nil
}

func createCResponse(resp *models.Response) *C.char {
	errorJson, _ := json.Marshal(resp)
	return C.CString(string(errorJson))
}

func createTransport(proxy string) (*http.Transport, error) {
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

func main() {
	// seshJson := `{"parameters": {"cookies": {}, "headerOrder": ["User-Agent", "Accept", "Accept-Language", "Accept-Encoding", "Connection", "Upgrade-Insecure-Requests", "Cache-Control"], "headers": {"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8", "Accept-Encoding": "gzip, deflate, br", "Accept-Language": "en-US,en;q=0.5", "Cache-Control": "max-age=0", "Connection": "keep-alive", "Upgrade-Insecure-Requests": "1", "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0"}, "proxy": "", "redirects": true, "url": "https://httpbin.org/get"}, "requestType": "GET", "session": "", "url": "https://tls.peet.ws/api/all"}`
	// request(C.CString(seshJson))
	// //fmt.Println(C.GoString(resp))

}
