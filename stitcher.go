package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const wordSeparator = ", "
const errorMsg = "ALL OTHER SOFTWARE TEAMS ARE USELESS, OUTSOURCING IS A SHAM"

// APIResult holds the result of a call to one of our wonderful external APIs
type APIResult struct {
	content string
	err     error
}

// Stitcher calls the Hello and World APIs to create a very useful string
func Stitcher(helloURL string, worldURL string) string {

	helloChannel := make(chan *APIResult, 1)
	worldChannel := make(chan *APIResult, 1)

	go getStringFromAPI(helloChannel, helloURL)
	go getStringFromAPI(worldChannel, worldURL)

	helloResult := <-helloChannel
	worldResult := <-worldChannel

	if helloResult.err != nil || worldResult.err != nil {
		return errorMsg
	}

	return helloResult.content + wordSeparator + worldResult.content
}

func getStringFromAPI(ch chan<- *APIResult, url string) {
	response, err := http.Get(url)

	if err != nil {
		ch <- &APIResult{"", err}
	} else {
		defer response.Body.Close()
		content, err := ioutil.ReadAll(response.Body)

		if err != nil {
			ch <- &APIResult{"", err}
		} else if response.StatusCode != http.StatusOK {
			ch <- &APIResult{"", errors.New("Non 200 response from API")}
		} else {
			ch <- &APIResult{string(content), nil}
		}
	}
}

func main() {
	fmt.Println(Stitcher("foo", "url"))
}
