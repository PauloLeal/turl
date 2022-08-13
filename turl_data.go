package turl

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type TURLData struct {
	Request struct {
		Method  string            `json:"method"`
		Url     string            `json:"url"`
		Headers map[string]string `json:"headers"`
		Query   map[string]string `json:"query"`
		Payload any               `json:"payload"`
	} `json:"request"`
	Response struct {
		StatusCode int               `json:"status_code"`
		Headers    map[string]string `json:"headers"`
		Payload    any               `json:"payload"`
	} `json:"response"`
	Error error
}

func (ld *TURLData) PrintText() {
	s := make([]string, 0)
	s = append(s, fmt.Sprintf("Request.Method: %s", ld.Request.Method))
	s = append(s, fmt.Sprintf("Request.Url: %s", ld.Request.Url))

	for k, v := range ld.Request.Headers {
		s = append(s, fmt.Sprintf("Request.Headers.%s: %s", k, v))
	}

	for k, v := range ld.Request.Query {
		s = append(s, fmt.Sprintf("Request.Query.%s: %s", k, v))
	}

	b, _ := json.Marshal(ld.Request.Payload)
	s = append(s, "Request.Payload: %s", string(b))

	s = append(s, fmt.Sprintf("Response.StatusCode: %d", ld.Response.StatusCode))

	for k, v := range ld.Response.Headers {
		s = append(s, fmt.Sprintf("Response.Headers.%s: %s", k, v))
	}

	b, _ = json.Marshal(ld.Response.Payload)
	s = append(s, "Response.Payload: %s", string(b))

	lockPrint(strings.Join(s, " | "))
}

func (ld *TURLData) PrintPretty() {
	s := fmt.Sprintf("Request  -- Method: %s\t\tURL: %s\n", ld.Request.Method, ld.Request.Url)

	s = fmt.Sprintf("%sRequest  -- Headers: \n", s)

	for k, v := range ld.Request.Headers {
		s = fmt.Sprintf("%sRequest  --\t%s = %s\n", s, k, v)
	}

	s = fmt.Sprintf("%sRequest  -- Query: \n", s)

	for k, v := range ld.Request.Query {
		s = fmt.Sprintf("%sRequest  --\t%s = %s\n", s, k, v)
	}

	s = fmt.Sprintf("%sRequest  -- Payload: \n", s)

	b, _ := json.MarshalIndent(ld.Request.Payload, "Request  -- \t", "    ")
	s = fmt.Sprintf("%sRequest  --\t%s\n", s, string(b))

	s = fmt.Sprintf("%s%s\n", s, strings.Repeat("-", 20))

	s = fmt.Sprintf("%sResponse -- Status Code: %d\n", s, ld.Response.StatusCode)
	s = fmt.Sprintf("%sResponse -- Headers: \n", s)
	for k, v := range ld.Response.Headers {
		s = fmt.Sprintf("%sResponse --\t%s = %s\n", s, k, v)
	}

	s = fmt.Sprintf("%sResponse -- Payload: \n", s)

	b, _ = json.MarshalIndent(ld.Response.Payload, "Response -- \t", "    ")
	s = fmt.Sprintf("%sResponse --\t%s\n", s, string(b))

	errs := "-"
	if ld.Error != nil {
		errs = ld.Error.Error()
	}
	s = fmt.Sprintf("%sResponse -- Error: %s\n", s, errs)
	s = fmt.Sprintf("%s%s", s, strings.Repeat("=", 120))

	lockPrint(s)
}

func (ld *TURLData) PrintJson() {
	b, _ := json.Marshal(ld)

	lockPrint(string(b))
}

var printMutex sync.Mutex

func lockPrint(s string) {
	printMutex.Lock()
	fmt.Println(s)
	printMutex.Unlock()
}
