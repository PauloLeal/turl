package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"gopkg.in/alecthomas/kingpin.v2"
)

type replacer func() string

var (
	method     = kingpin.Flag("method", "Http method").Short('X').Required().String()
	targetUrl  = kingpin.Arg("url", "Target url").Required().String()
	count      = kingpin.Flag("count", "Number of requets").Short('c').Default("1").Int()
	headers    = kingpin.Flag("header", "Http headers to send").Short('H').StringMap()
	pathParams = kingpin.Flag("path", "Path parameters to send").Short('p').StringMap()
	payload    = kingpin.Flag("data-raw", "Request body to send").String()
	outFormat  = kingpin.Flag("output-format", "Output format to report").Default("pretty").Enum("json", "text", "pretty")
	errorsOnly = kingpin.Flag("errors-only", "Show only errors (not http status errors)").Bool()

	replacements = map[string]replacer{
		"{{UUID}}": func() string { return uuid.New().String() },
		"{{XULID}}": func() string {
			hexValue := fmt.Sprintf("%x", string(ulid.Make().Bytes()))
			hexBytes, _ := hex.DecodeString(hexValue)
			u, _ := uuid.FromBytes(hexBytes)

			return u.String()
		},
		"{{ULID}}":      func() string { return ulid.Make().String() },
		"{{RANDINT-3}}": func() string { return fmt.Sprintf("%03d", randN(3)) },
		"{{RANDINT-5}}": func() string { return fmt.Sprintf("%05d", randN(5)) },
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randN(n int) int {
	d := math.Pow(10, float64(n))
	return rand.Intn(int(d))
}

func main() {
	kingpin.Parse()

	if *count < 0 {
		log.Fatal("count cannot be < 0")
	}

	var wg sync.WaitGroup
	for i := 0; i < *count; i++ {
		rMethod := makeReplacements(*method)
		rUrl := makeReplacements(*targetUrl)

		rHeaders := make(map[string]string)
		for k, v := range *headers {
			rHeaders[k] = makeReplacements(v)
		}

		rBody := makeReplacements(*payload)

		for k, v := range *pathParams {
			rUrl = strings.ReplaceAll(rUrl, fmt.Sprintf(":%s", k), v)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			ld := HttpLogData{}

			ld.Request.Method = rMethod
			ld.Request.Url = rUrl
			ld.Request.Headers = rHeaders
			json.Unmarshal([]byte(rBody), &ld.Request.Payload)

			status, body, headers, err := makeRequest(rMethod, rUrl, rHeaders, []byte(rBody))

			ld.Response.StatusCode = status
			json.Unmarshal(body, &ld.Response.Payload)
			ld.Response.Headers = headers

			ld.Error = err

			if *errorsOnly && ld.Error == nil {
				return
			}

			switch *outFormat {
			case "pretty":
				ld.PrintPretty()
			case "json":
				ld.PrintJson()
			case "text":
				ld.PrintText()
			}
		}()
	}

	wg.Wait()
}

func makeReplacements(s string) string {
	for k, v := range replacements {
		for {
			oldS := s
			s = strings.Replace(s, k, v(), 1)

			if oldS == s {
				break
			}
		}
	}

	return s
}

func makeRequest(method string, url string, headers map[string]string, payload []byte) (int, []byte, map[string]string, error) {
	var res *resty.Response
	var err error
	client := resty.New()

	req := client.R().
		SetHeaders(headers).
		SetBody(payload)

	res, err = req.Execute(strings.ToUpper(method), url)

	responseHeaders := make(map[string]string)
	for k, v := range res.Header() {
		responseHeaders[k] = strings.Join(v, ", ")
	}

	return res.StatusCode(), res.Body(), responseHeaders, err
}
