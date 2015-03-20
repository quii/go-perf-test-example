// Package stitcher is a simple of example of using Go's tools to produce well tested, performant software
package stitcher

import (
	"errors"
	"io/ioutil"
	"net/http"
)

const wordSeparator = ", "
const errorMsg = "ALL OTHER SOFTWARE TEAMS ARE USELESS, OUTSOURCING IS A SHAM"

type apiResult struct {
	content string
	err     error
}

// Stitcher calls the Hello and World APIs to create a very useful string
func Stitcher(helloURL string, worldURL string) string {

	helloChannel := make(chan *apiResult, 1)
	worldChannel := make(chan *apiResult, 1)

	go getStringFromAPI(helloChannel, helloURL)
	go getStringFromAPI(worldChannel, worldURL)

	helloResult := <-helloChannel
	worldResult := <-worldChannel

	if helloResult.err != nil || worldResult.err != nil {
		return errorMsg
	}

	return helloResult.content + wordSeparator + worldResult.content
}

func getStringFromAPI(ch chan<- *apiResult, url string) {
	response, err := http.Get(url)

	if err != nil {
		ch <- &apiResult{"", err}
	} else {
		defer response.Body.Close()
		content, err := ioutil.ReadAll(response.Body)

		if err != nil {
			ch <- &apiResult{"", err}
		} else if response.StatusCode != http.StatusOK {
			ch <- &apiResult{"", errors.New("Non 200 response from API")}
		} else {
			ch <- &apiResult{string(content), nil}
		}
	}
}
