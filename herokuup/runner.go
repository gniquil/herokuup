package herokuup

import (
	"fmt"
	"net/http"
)

type Config struct {
	Urls       []string `json:"urls"`
	From       string   `json:"from"`
	Tos        []string `json:"tos"`
	Sendonfail bool     `json:"sendonfail"`
	Serveraddr string   `json:"serveraddr"`
}

type Response struct {
	url    string
	status int
}

func checkUrl(url string, resChan chan<- Response) {
	res, err := http.Get(url)

	var status int // zero'd to 0
	if err == nil {
		status = res.StatusCode
	}

	fmt.Printf("%s -> %d\n", url, status)

	resChan <- Response{url: url, status: status}
}

func Run(path string) {
	fmt.Println("Loading config...")

	config := loadConfig(path)

	fmt.Println("Setting up...")

	var responses []Response

	resChan := make(chan Response)

	fmt.Println("Sending out gophers...")

	for _, url := range config.Urls {
		go checkUrl(url, resChan)
	}

	fmt.Println("Listening for responses...")

	for len(responses) < len(config.Urls) {
		res := <-resChan
		responses = append(responses, res)
	}

	fmt.Println("Generating messages...")

	sendMessage(config, responses)
}
