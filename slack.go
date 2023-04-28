package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

const webhook = "https://hooks.slack.com/services/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"

// PostWebhook makes a POST request containing message to the specfied url. url is webhook that setup on a slack application.
func postWebhook(message string) error {
	type body struct {
		Text string `json:"text"`
	}
	payload := &body{Text: message}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}
	_, err = http.Post(webhook, "application/json", &buf)
	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}
	return nil
}
