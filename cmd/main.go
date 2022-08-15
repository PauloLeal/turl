package main

import (
	"log"
	"sync"
	"time"

	"github.com/PauloLeal/turl"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	targetUrl      = kingpin.Arg("url", "Target url").Required().String()
	followLocation = kingpin.Flag("location", "Follow location redirects").Short('L').Bool()
	method         = kingpin.Flag("request", "Http method").Short('X').Required().String()
	count          = kingpin.Flag("count", "Number of requets").Short('c').Default("1").Int()
	delay          = kingpin.Flag("delay", "Time to wait before starting").Default("0s").Duration()
	loop           = kingpin.Flag("loop", "Repeat `count` requests for `loop` times").Default("1").Int()
	loopInterval   = kingpin.Flag("loop-interval", "Time between loop iteractions").Default("5s").Duration()
	headers        = kingpin.Flag("header", "Http headers to send").Short('H').StringMap()
	payload        = kingpin.Flag("data-raw", "Request body to send").Short('d').String()
	pathParams     = kingpin.Flag("path", "Path parameters to send").Short('p').StringMap()
	queryParams    = kingpin.Flag("query", "Query parameters to send").Short('Q').StringMap()
	outFormat      = kingpin.Flag("output-format", "Output format to report").Default("pretty").Enum("json", "text", "pretty")
	noOutput       = kingpin.Flag("no-output", "Supress requests output").Bool()
	statusOnly     = kingpin.Flag("status-only", "Show only requests with this status").Int()
	errorsOnly     = kingpin.Flag("errors-only", "Show only errors (not http status errors)").Bool()
)

func main() {
	kingpin.Parse()

	if *count < 0 {
		log.Fatal("count cannot be < 0")
	}

	if *loop < 1 {
		*loop = 1
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

	time.Sleep(*delay)

	var wg sync.WaitGroup
	for l := 0; l < *loop; l++ {
		if l > 0 {
			time.Sleep(*loopInterval)
		}
		for c := 0; c < *count; c++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				ld := turl.MakeRequest(params)

				if *noOutput {
					return
				}

				if *errorsOnly && ld.Error == nil {
					return
				}

				if *statusOnly > 0 && ld.Response.StatusCode != *statusOnly {
					return
				}

				switch *outFormat {
				case "json":
					ld.PrintJson()
				case "text":
					ld.PrintText()
				default:
					ld.PrintPretty()
				}
			}()
		}

		wg.Wait()
	}
}
