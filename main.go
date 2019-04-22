package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	concurrency := flag.Int("concurrency", 10, "Number of parallel threads to send the request from")
	numberOfRequests := flag.Int("number", 100, "Number of requests to send")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		panic("URL must be supplied")
	}

	url := args[0]

	log.Printf(
		"Starting job with concurrency = %d, numberOfRequests = %d and url = %s\n",
		*concurrency,
		*numberOfRequests,
		url,
	)

	requestChannel := make(chan int)
	responseChannel := make(chan Response)
	resultChannel := make(chan Result)

	for i := 0; i < *concurrency; i++ {
		go sendingThread(url, requestChannel, responseChannel)
	}

	go aggregatingThread(responseChannel, resultChannel, *numberOfRequests)

	for j := 0; j < *numberOfRequests; j++ {
		requestChannel <- j
	}
	close(requestChannel)

	result := <-resultChannel
	close(responseChannel)

	log.Printf("Total Requests: %d\n", result.Total)
	log.Printf("Total Passed  : %d\n", result.TotalPass)
	log.Printf("Total Failed  : %d\n", result.Total-result.TotalPass)
	log.Printf("Pass Percent  : %.2f%%\n", result.passPercentage())
}

func sendingThread(url string, requestChannel <-chan int, responseChannel chan<- Response) {
	for range requestChannel {
		res, err := http.Get(url)

		pass := false
		statusCode := 0

		if err != nil {
			log.Printf("Error sending request %s", err.Error())
		} else {
			if res.StatusCode >= 200 && res.StatusCode < 300 {
				pass = true
			}

			statusCode = res.StatusCode
		}

		responseChannel <- Response{
			Code: statusCode,
			Pass: pass,
		}
	}
}

func aggregatingThread(responseChannel <-chan Response, resultChannel chan<- Result, expectedRequests int) {

	result := Result{
		TotalPass: 0,
		Total:     0,
		Expected:  expectedRequests,
	}

	for res := range responseChannel {
		result.Total++

		if res.Pass {
			result.TotalPass = result.TotalPass + 1
		}

		if result.donePercentage()%10 == 0 {
			log.Printf(
				"Progress: Requests sent = %d, Successful = %d, Failed = %d\n",
				result.Total,
				result.TotalPass,
				result.Total-result.TotalPass,
			)
		}

		if result.Total >= expectedRequests {
			resultChannel <- result
			close(resultChannel)
		}
	}
}

type Response struct {
	Code int
	Pass bool
}

type Result struct {
	TotalPass int
	Total     int
	Expected  int
}

func (r Result) passPercentage() float32 {
	return (float32(r.TotalPass) / float32(r.Total)) * 100
}

func (r Result) donePercentage() int {
	return int((float32(r.Total) / float32(r.Expected)) * 100)
}
