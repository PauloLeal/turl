package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

var printMutex sync.Mutex

type httpRequestLogData struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Payload any               `json:"payload"`
}

type httpResponseLogData struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Payload    any               `json:"payload"`
}

type HttpLogData struct {
	Request  httpRequestLogData  `json:"request"`
	Response httpResponseLogData `json:"response"`
}

func (ld *HttpLogData) PrintText() {
	s := make([]string, 0)
	s = append(s, fmt.Sprintf("Request.Method: %s", ld.Request.Method))
	s = append(s, fmt.Sprintf("Request.Url: %s", ld.Request.Url))

	for k, v := range ld.Request.Headers {
		s = append(s, fmt.Sprintf("Request.Headers.%s: %s", k, v))
	}

	b, _ := json.Marshal(ld.Request.Payload)
	s = append(s, "Request.Payload: %s", string(b))

	s = append(s, fmt.Sprintf("Response.StatusCode: %d", ld.Response.StatusCode))

	for k, v := range ld.Response.Headers {
		s = append(s, fmt.Sprintf("Response.Headers.%s: %s", k, v))
	}

	b, _ = json.Marshal(ld.Response.Payload)
	s = append(s, "Response.Payload: %s", string(b))

	printMutex.Lock()
	fmt.Println(strings.Join(s, " | "))
	printMutex.Unlock()
}

func (ld *HttpLogData) PrintPretty() {
	s := fmt.Sprintf("Request  -- Method: %s\t\tURL: %s\n", ld.Request.Method, ld.Request.Url)
	s = fmt.Sprintf("%sRequest  -- Headers: \n", s)

	for k, v := range ld.Request.Headers {
		s = fmt.Sprintf("%sRequest  --\t%s = %s\n", s, k, v)
	}

	s = fmt.Sprintf("%sRequest  -- Payload: \n", s)

	b, _ := json.MarshalIndent(ld.Request.Payload, "Request  -- \t", "    ")
	s = fmt.Sprintf("%sRequest  --\t%s\n", s, string(b))

	s = fmt.Sprintf("%s\n%s\n", s, strings.Repeat("-", 20))

	s = fmt.Sprintf("%sResponse -- Status Code: %d\n", s, ld.Response.StatusCode)
	s = fmt.Sprintf("%sResponse -- Headers: \n", s)
	for k, v := range ld.Response.Headers {
		s = fmt.Sprintf("%sResponse --\t%s = %s\n", s, k, v)
	}

	s = fmt.Sprintf("%sResponse -- Payload: \n", s)

	b, _ = json.MarshalIndent(ld.Response.Payload, "Response -- \t", "    ")
	s = fmt.Sprintf("%sResponse --\t%s\n", s, string(b))

	s = fmt.Sprintf("%s\n%s", s, strings.Repeat("=", 120))

	printMutex.Lock()
	fmt.Println(s)
	printMutex.Unlock()
}

func (ld *HttpLogData) PrintJson() {
	b, _ := json.Marshal(ld)

	printMutex.Lock()
	fmt.Println(string(b))
	printMutex.Unlock()
}
