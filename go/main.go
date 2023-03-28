package main

import (
	"C"
	"bytes"
	"chromeRequests/models"
	"encoding/json"
	"net/url"

	http "github.com/saucesteals/fhttp"
	"github.com/saucesteals/mimic"
)
import (
	"fmt"
	"io"
	"sync"
)

var Sessions = make(map[string]models.Session)
var latestVersion = mimic.MustGetLatestVersion(mimic.PlatformWindows)
var m, _ = mimic.Chromium(mimic.BrandChrome, latestVersion)
var cachedTransport = map[string]*LockableTransport{}

func createTransport(proxy string) (*LockableTransport, error) {
	var LockableTransport = &LockableTransport{}

	if trans, ok := cachedTransport[proxy]; ok {
		return trans, nil
	}

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)

		if err != nil {
			return nil, err
		}

		transport := m.ConfigureTransport(&http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		})

		LockableTransport.Transport = transport
		LockableTransport.Lock = &sync.Mutex{}

		cachedTransport[proxy] = LockableTransport

		return cachedTransport[proxy], nil
	}

	transport := m.ConfigureTransport(&http.Transport{})
	LockableTransport.Transport = transport
	LockableTransport.Lock = &sync.Mutex{}

	return LockableTransport, nil
}

//export request
func request(cParams *C.char) *C.char {
	var client http.Client = http.Client{}

	var req *http.Request

	params := C.GoString(cParams)

	data := models.SessionParameters{}

	err := json.Unmarshal([]byte(params), &data)

	if err != nil {
		return createCResponse(&models.Response{Error: err.Error()})
	}

	req, err = http.NewRequest(data.RequestType, data.Parameters.URL, nil)

	if err != nil {
		return createCResponse(&models.Response{Error: err.Error()})
	}

	if data.Parameters.Json != "" {
		req.Body = io.NopCloser(bytes.NewBuffer([]byte(data.Parameters.Json)))
	}

	if data.Parameters.Form != nil {
		formData := url.Values{}
		for key, value := range data.Parameters.Form {
			formData.Add(key, value)
		}
		req.Body = io.NopCloser(bytes.NewBufferString(formData.Encode()))
	}

	transport, err := createTransport(data.Parameters.Proxy)

	if err != nil {
		return createCResponse(&models.Response{Error: err.Error()})
	}

	transport.Lock.Lock()
	defer transport.Lock.Unlock()
	client.Transport = transport.Transport

	if data.Parameters.Headers != nil {
		for k, v := range data.Parameters.Headers {
			req.Header[k] = []string{v}
		}
		req.Header[http.HeaderOrderKey] = data.Parameters.HeaderOrder
		req.Header[http.PHeaderOrderKey] = m.PseudoHeaderOrder()
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

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return createCResponse(&models.Response{Error: err.Error()})
	}

	headersMap := map[string]string{}

	cookieMap := map[string]string{}

	for _, cookie := range resp.Cookies() {
		cookieMap[cookie.Name] = cookie.Value
	}

	client.CheckRedirect = nil

	return createCResponse(&models.Response{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Cookies:    cookieMap,
		Headers:    headersMap,
		Url:        resp.Request.URL.String(),
	})
}

func createCResponse(resp *models.Response) *C.char {
	errorJson, _ := json.Marshal(resp)
	return C.CString(string(errorJson))
}

func main() {
	// RequestStructure := models.SessionParameters{
	// 	Parameters: &models.RequestParameters{
	// 		URL: "https://tls.peet.ws/api/all",
	// 		Headers: map[string]string{
	// 			"user-agent":      "Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0",
	// 			"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	// 			"accept-language": "en-US,en;q=0.5",
	// 			"accept-encoding": "gzip, deflate, br",
	// 		},
	// 		HeaderOrder: []string{
	// 			"user-agent",
	// 			"accept",
	// 			"accept-language",
	// 			"accept-encoding",
	// 		},
	// 		Cookies: map[string]string{},
	// 	},
	// 	RequestType: "GET",
	// }
	// seshJson, _ := json.Marshal(RequestStructure)
	// log.Println(string(seshJson))
	resp := request(C.CString(string(`{"parameters": {"cookies": {}, "form": {}, "headerOrder": [], "headers": {}, "json": "{}", "proxy": "", "redirects": true, "url": "https://httpbin.org/post"}, "requestType": "POST", "url": "https://httpbin.org/post"}`)))
	fmt.Println(C.GoString(resp))

}

type LockableTransport struct {
	Transport *http.Transport
	Lock      *sync.Mutex
}

//make a command to build a c-shared library
//go:generate go build -buildmode=c-shared -o chromeRequests.so main.go
