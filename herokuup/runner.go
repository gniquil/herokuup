package herokuup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
)

type RunnerConfig struct {
	Urls       []string `json:"checks"`
	From       string   `json:"from_email"`
	Tos        []string `json:"to_emails"`
	Sendonfail bool     `json:"send_only_on_fail"`
	Serveraddr string   `json:"server_address"`
}

type Runner struct {
	config      *RunnerConfig
	responses   []Response
	totalFailed int
}

func NewRunner(path string) *Runner {
	runner := new(Runner)
	runner.parseConfig(path)
	return runner
}

func (runner *Runner) Run() {
	fmt.Println(runner.config)

	resChan := make(chan Response)
	// allResponses := make([]Response, 0)

	for _, url := range runner.config.Urls {
		go runner.checkUrl(url, resChan)
	}

	for {
		if len(runner.responses) == len(runner.config.Urls) {
			break
		}

		select {
		case res := <-resChan:
			runner.responses = append(runner.responses, res)
			if res.failed() {
				runner.totalFailed++
			}
		}
	}

	runner.sendMessage(runner.craftMessage())
}

func (runner *Runner) parseConfig(path string) {
	config := new(RunnerConfig)

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		panic(err)
	}

	runner.config = config
}

func (runner *Runner) checkUrl(url string, resChan chan<- Response) {
	res, err := http.Get(url)
	fmt.Println(url, res.StatusCode)
	if err != nil {
		resChan <- Response{url: url, status: 0}
	} else {
		resChan <- Response{url: url, status: res.StatusCode}
	}
}

func (runner *Runner) craftMessage() string {
	message := ""
	for _, response := range runner.responses {
		if response.failed() {
			message += fmt.Sprintf("%s returned status code %d\n", response.url, response.status)
		}
	}
	var finalMessage string
	if runner.totalFailed > 0 {
		finalMessage = fmt.Sprintf("Subject: [herokuup] %d urls are down!\n\n", runner.totalFailed)
		finalMessage += message + "\n"
		finalMessage += fmt.Sprintf("%d out of %d failed.", runner.totalFailed, len(runner.responses))
	} else {
		finalMessage = "Subject: [herokuup] All urls are up!\n\n"
		finalMessage += fmt.Sprintf("%d out of %d passed.", len(runner.responses), len(runner.responses))
	}
	finalMessage += "\n"
	return finalMessage
}

func (runner *Runner) sendMessage(message string) {
	if runner.config.Sendonfail && runner.totalFailed == 0 {
		fmt.Println(message)
		return
	}

	err := smtp.SendMail(
		runner.config.Serveraddr,
		nil,
		runner.config.From,
		runner.config.Tos,
		[]byte(message),
	)
	if err != nil {
		panic(err)
	}
}
