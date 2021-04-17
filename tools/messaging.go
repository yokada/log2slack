package tools

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type MessageBody struct {
	Text string `json:"text"`
}

func SendAsync(b []byte) (<-chan string, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	message := &MessageBody{Text: string(b)}

	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	log.Printf("webhookUrl=%v", config.SlackAPI.Messaging.WebhookURL)
	req, err := http.NewRequest("POST", config.SlackAPI.Messaging.WebhookURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}

	done := make(chan string)

	go func() {
		defer close(done)

		res, err := client.Do(req)
		if err != nil {
			done <- err.Error()
			return
		}

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			done <- err.Error()
			return
		}

		done <- string(resBody)
	}()

	return done, nil
}
