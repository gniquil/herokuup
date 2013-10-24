package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
)

type Config struct {
	Checks         []string `json:"checks"`
	FromEmail      string   `json:"from_email"`
	ToEmails       []string `json:"to_emails"`
	SendOnlyOnFail bool     `json:"send_only_on_fail"`
	ServerAddress  string   `json:"server_address"`
}

func main() {
	config := ParseConfig(os.Args[1])

	channel := make(chan map[string]int)

	for _, url := range config.Checks {
		go Check(url, channel)
	}

	responses := []map[string]int{}

	for i := 0; i < len(config.Checks); i++ {
		checkResponse := <-channel
		responses = append(responses, checkResponse)
	}

	var totalFailed int
	message, totalFailed := CraftMessage(responses)

	if config.SendOnlyOnFail {
		if totalFailed > 0 {
			SendMessage(message, config)
		} else {
			fmt.Println(message)
		}
	} else {
		SendMessage(message, config)
	}
}

func SendMessage(message string, config Config) {
	err := smtp.SendMail(
		config.ServerAddress,
		nil,
		config.FromEmail,
		config.ToEmails,
		[]byte(message),
	)
	if err != nil {
		panic(err)
	}
}

func CraftMessage(responses []map[string]int) (string, int) {
	message := ""
	totalFailed := 0
	for _, response := range responses {
		for key, val := range response {
			if val != 200 {
				message += fmt.Sprintf("%s returned status code %d\n", key, val)
				totalFailed++
			}
		}
	}
	var finalMessage string
	if totalFailed > 0 {
		finalMessage = fmt.Sprintf("Subject: [herokuup] %d urls are down!\n\n", totalFailed)
		finalMessage += message + "\n"
		finalMessage += fmt.Sprintf("%d out of %d failed.", totalFailed, len(responses))
	} else {
		finalMessage = "Subject: [herokuup] All urls are up!\n\n"
		finalMessage += fmt.Sprintf("%d out of %d passed.", len(responses), len(responses))
	}
	finalMessage += "\n"
	return finalMessage, totalFailed
}

func Check(url string, channel chan<- map[string]int) {
	res, err := http.Get(url)
	if err != nil {
		channel <- map[string]int{url: 0}
	} else {
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
