package turl

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/PauloLeal/turl/integration"
)

type TURL struct {
	Method         string
	URL            string
	Headers        map[string]string
	PathParams     map[string]string
	QueryParams    map[string]string
	Payload        string
	FollowLocation bool
}

type FHttpRequester func(method string, url string, headers map[string]string, query map[string]string, payload string) (int, string, map[string]string, error)

var HttpRequester FHttpRequester

func init() {
	HttpRequester = integration.MakeHttpRequest
}

func MakeRequest(t TURL) TURLData {
	ld := TURLData{}

	rMethod := makeReplacements(t.Method)

	rUrl := makeReplacements(t.URL)

	rHeaders := make(map[string]string)
	for k, v := range t.Headers {
		rHeaders[k] = makeReplacements(v)
	}

	rQuery := make(map[string]string)

	purl, err := url.Parse(rUrl)
	if err != nil {
		ld.Error = err
		return ld
	}

	for k, vl := range purl.Query() {
		rQuery[k] = makeReplacements(strings.Join(vl, ", "))
	}

	for k, v := range t.QueryParams {
		rQuery[k] = makeReplacements(v)
	}

	ld.Request.Query = rQuery

	rBody := makeReplacements(t.Payload)

	for k, v := range t.PathParams {
		rUrl = strings.ReplaceAll(rUrl, fmt.Sprintf(":%s", k), v)
	}

	ld.Request.Method = rMethod
	ld.Request.Url = rUrl
	ld.Request.Headers = rHeaders
	json.Unmarshal([]byte(rBody), &ld.Request.Payload)

	status, body, headers, err := HttpRequester(rMethod, rUrl, rHeaders, rQuery, rBody)

	ld.Response.StatusCode = status
	json.Unmarshal([]byte(body), &ld.Response.Payload)
	ld.Response.Headers = headers

	ld.Error = err

	return ld
}
