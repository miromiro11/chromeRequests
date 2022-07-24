package main

// #include <stdio.h>
// #include <stdlib.h>
//
// static void myprint(char* s) {
//   printf("%s\n", s);
// }

import (
	"C"
	"encoding/json"
	"net/url"

	http "github.com/saucesteals/fhttp"
	"github.com/saucesteals/mimic"
)

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
)

var Sessions =  make(map[string]Session)
var userAgent = fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", m.Version())
var latestVersion = mimic.MustGetLatestVersion(mimic.PlatformWindows)
var m, _ = mimic.Chromium(mimic.BrandChrome, latestVersion)
var globalHeaders = make(map[string]string)
var cleanTransport = &http.Transport{}

type Session struct {
	Client *http.Client
	Headers map[string]string
}

type Response struct {
	StatusCode int
	Body       string
	Cookies   map[string]string
	Headers   map[string]string
}

type RequestParameters struct {
	URL     string   `json:"url"`
	Proxy   string   `json:"proxy"`
	Headers map[string]string `json:"headers"`
	Form    map[string]string   `json:"FORM"`
	Json   string   `json:"JSON"`
}

type sessionParamters struct{
	Session string `json:"session"`
	RequestType string `json:"requestType"`
	Paramters RequestParameters `json:"paramters"`
	Proxy string `json:"proxy"`
}

type headerChange struct{
	Session string `json:"session"`
	Headers map[string]string `json:"headers"`
}

func createTransport(proxy string) *http.Transport {
	if len(proxy) != 0 {
		proxyUrl, err := url.Parse(proxy)
		check(err)
		return &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	} else {
		return &http.Transport{}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
//export changeProxy
func changeProxy(params *C.char){
	var sessionParameters sessionParamters
	json.Unmarshal([]byte(C.GoString(params)), &sessionParameters)
	proxy := sessionParameters.Proxy
	session := sessionParameters.Session
	Sessions[session].Client.Transport = m.ConfigureTransport(createTransport(proxy))
}

//export addHeaders 
func addHeaders(Params *C.char) {
	var sessionParameters headerChange
	json.Unmarshal([]byte(C.GoString(Params)), &sessionParameters)
	session := sessionParameters.Session
	headers := sessionParameters.Headers
	for key, value := range headers {
		Sessions[session].Headers[key] = value
	}
}

//export removeHeaders 
func removeHeaders(Params *C.char) {
	var sessionParameters headerChange
	json.Unmarshal([]byte(C.GoString(Params)), &sessionParameters)
	session := sessionParameters.Session
	headers := sessionParameters.Headers
	for key, _ := range headers {
		delete(Sessions[session].Headers, key)
	}
}

//export createSession
func createSession(proxy *C.char) *C.char {
	proxy_ := C.GoString(proxy)
	newUUID, err := exec.Command("uuidgen").Output()
	newUUID = newUUID[:len(newUUID)-1]
	check(err)
	Sessions[string(newUUID)] = Session{
		Client:  &http.Client{Transport: m.ConfigureTransport(createTransport(proxy_))},
		Headers: make(map[string]string),
	}
	return C.CString(string(newUUID))
}

//export getMap
func getMap() *C.char{
	json, _ := json.Marshal(Sessions)
	return C.CString(string(json))
}

//export request
func request(params *C.char) *C.char {
	var client *http.Client
	var req *http.Request
	var err error
	params_ := C.GoString(params)
	data := sessionParamters{}
	json.Unmarshal([]byte(params_), &data)
	if data.Session != "" {
		client = Sessions[data.Session].Client
	} else {
		client = &http.Client{
			Transport: m.ConfigureTransport(createTransport("")),
		}
	}
	if data.RequestType == "GET" {
		req,err = http.NewRequest("GET", data.Paramters.URL, nil)
		check(err)
	} else if data.RequestType == "POST" {
		req,err = http.NewRequest("POST", data.Paramters.URL, nil)
		check(err)
		if len(data.Paramters.Form) != 0 {
			url := url.Values{}
			for key, value := range data.Paramters.Form {
				url.Add(key, value)
			}
			req,_ = http.NewRequest("POST", data.Paramters.URL, bytes.NewBufferString(url.Encode()))
		} else if data.Paramters.Json != "" {
			req,_ = http.NewRequest("POST", data.Paramters.URL, bytes.NewBuffer([]byte(data.Paramters.Json)))
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
	if data.Paramters.Proxy != "" {
		client.Transport = m.ConfigureTransport(createTransport(data.Paramters.Proxy))
	}
	if data.Paramters.Headers != nil {
		for k, v := range data.Paramters.Headers {
			req.Header.Add(k, v)
		}
	}
	for k, v := range Sessions[data.Session].Headers {
		req.Header.Add(k, v)
	}
	if data.Paramters.Json != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if len(data.Paramters.Form) != 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	cookies := resp.Cookies()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	headersMap := make(map[string]string)
	for key, value := range resp.Header {
		headersMap[key] = value[0]
	}
	cookieMap := make(map[string]string)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}
	response := Response{resp.StatusCode, string(body), cookieMap, headersMap}
	json, _ := json.Marshal(response)
	check(err)
	client.Transport = cleanTransport
	return C.CString(string(json))
}

func main() {
	
}
