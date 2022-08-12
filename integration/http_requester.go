package integration

import (
	"strings"

	"github.com/go-resty/resty/v2"
)

func MakeHttpRequest(method string, url string, headers map[string]string, query map[string]string, payload string) (int, string, map[string]string, error) {
	var res *resty.Response
	var err error
	client := resty.New()

	req := client.R().
		SetHeaders(headers).
		SetQueryParams(query).
		SetBody(payload)

	res, err = req.Execute(strings.ToUpper(method), url)

	responseHeaders := make(map[string]string)
	for k, v := range res.Header() {
		responseHeaders[k] = strings.Join(v, ", ")
	}

	return res.StatusCode(), string(res.Body()), responseHeaders, err
}
