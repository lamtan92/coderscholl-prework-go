package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type reponseInfo struct {
	status   int
	bytes    int64
	duration time.Duration
}

type summaryInfo struct {
	requested int64
	responded int64
}

func main() {
	requests := flag.Int64("n", 1, "Number of requests to perform")
	concurrency := flag.Int64("c", 1, "Number of multiple requests to make at a time")
	fmt.Println(requests, concurrency)
	flag.Parse()
	if flag.NArg() == 0 || *requests == 0 || *requests < *concurrency {
		flag.PrintDefaults()
		os.Exit(0)
	}
	link := flag.Arg(0)
	c := make(chan reponseInfo)
	summary := summaryInfo{}
	for i := int64(0); i < *concurrency; i++ {
		summary.requested++
		go checkLink(link, c)
	}

	for response := range c {

		if summary.requested < *requests {
			summary.requested++
			go checkLink(link, c)
		}

		summary.responded++
		fmt.Println(response)
		if summary.responded == summary.requested {
			break
		}
	}
}

func checkLink(link string, c chan reponseInfo) {
	start := time.Now()
	res, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	read, _ := io.Copy(ioutil.Discard, res.Body)
	c <- reponseInfo{
		status:   res.StatusCode,
		bytes:    read,
		duration: time.Now().Sub(start),
	}
}
