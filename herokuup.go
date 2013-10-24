package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	Checks         []string `json:"checks"`
	FromEmail      string   `json:"from_email"`
	ToEmail        string   `json:"to_email"`
	SendOnlyOnFail bool     `json:"send_only_on_fail"`
}

func main() {
	config := ParseConfig(os.Args[1])

	channel := make(chan map[string]int)

	for _, url := range config.Checks {
		go Check(url, channel)
	}

	totalFailed := 0
	for i := 0; i < len(config.Checks); i++ {
		checkResponse := <-channel
		for _, val := range checkResponse {
			if val != 200 {
				totalFailed++
			}
		}
	}

	fmt.Printf("%d out of %d failed!", totalFailed, len(config.Checks))
}

func Check(url string, channel chan<- map[string]int) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		channel <- map[string]int{url: 0}
	} else {
		fmt.Println(url, res.StatusCode)
		channel <- map[string]int{url: res.StatusCode}
	}
}

func ParseConfig(configPath string) (config Config) {
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		panic(err)
	}

	return
}
