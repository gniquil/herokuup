package herokuup

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"regexp"
	"strings"
)

func loadConfig(path string) *Config {
	config := new(Config)

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		panic(err)
	}

	return config
}

func Align(text string) string {
	re := regexp.MustCompile(`^\s+`)
	result := ""
	stringArray := strings.Split(text, "\n")
	for i := 1; i < len(stringArray)-1; i++ {
		result += re.ReplaceAllString(stringArray[i], "") + "\n"
	}
	return result
}

func craftMessage(responses []Response) (string, int) {
	message := ""
	totalFailed := 0

	for _, response := range responses {
		if response.status != 200 {
			message += fmt.Sprintf("%s returned status code %d\n", response.url, response.status)
			totalFailed++
		}
	}

	if totalFailed > 0 {
		message = fmt.Sprintf(Align(`
      Subject: [herokuup] %d urls are down!

      %s
      %d out of %d failed.
    `), totalFailed, message, totalFailed, len(responses))
	} else {
		message = fmt.Sprintf(Align(`
      Subject: [herokuup] All urls are up!

      %d out of %d passed.
    `), len(responses), len(responses))
	}

	return message, totalFailed
}

func sendMessage(config *Config, responses []Response) {
	message, totalFailed := craftMessage(responses)

	if config.Sendonfail && totalFailed == 0 {
		fmt.Print(message)
		return
	}

	err := smtp.SendMail(
		config.Serveraddr,
		nil,
		config.From,
		config.Tos,
		[]byte(message),
	)
	if err != nil {
		panic(err)
	}
}
