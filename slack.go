package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	webhook    string
	webhookOld string
)

func initSlack(webhookVar, webhookOldVar string) {
	webhook = webhookOldVar
	webhookOld = webhookOldVar
}

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
		_, err = http.Post(webhookOld, "application/json", &buf)
		if err != nil {
			logrus.Errorf("%v", err)
			return err
		}
	}
	return nil
}
