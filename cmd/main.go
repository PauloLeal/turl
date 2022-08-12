package main

import (
	"log"
	"sync"

	"github.com/PauloLeal/turl"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	targetUrl      = kingpin.Arg("url", "Target url").Required().String()
	followLocation = kingpin.Flag("location", "Follow location redirects").Short('L').Bool()
	method         = kingpin.Flag("request", "Http method").Short('X').Required().String()
	count          = kingpin.Flag("count", "Number of requets").Short('c').Default("1").Int()
	headers        = kingpin.Flag("header", "Http headers to send").Short('H').StringMap()
	payload        = kingpin.Flag("data-raw", "Request body to send").Short('d').String()
	pathParams     = kingpin.Flag("path", "Path parameters to send").Short('p').StringMap()
	queryParams    = kingpin.Flag("query", "Query parameters to send").Short('Q').StringMap()
	outFormat      = kingpin.Flag("output-format", "Output format to report").Default("pretty").Enum("json", "text", "pretty")
	errorsOnly     = kingpin.Flag("errors-only", "Show only errors (not http status errors)").Bool()
)

func main() {
	kingpin.Parse()

	if *count < 0 {
		log.Fatal("count cannot be < 0")
	}

	params := turl.TURL{
		Method:         *method,
		URL:            *targetUrl,
		Headers:        *headers,
		PathParams:     *pathParams,
		QueryParams:    *queryParams,
		Payload:        *payload,
		FollowLocation: *followLocation,
	}

	var wg sync.WaitGroup
	for i := 0; i < *count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ld := turl.MakeRequest(params)
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
