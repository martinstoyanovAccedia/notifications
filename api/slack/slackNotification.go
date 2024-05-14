package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"notifications/api/pooling"
	"notifications/types"
)

func Webhook(msg *SlackRequestBody) {
	cfg := types.SlackMsgEnvConfig()
	webhookUrl := cfg.SlackReportWebhook

	err := SendSlackNotification(webhookUrl, msg)
	if err != nil {
		log.Print("unable to send message to slack")
	}
}

func SendSlackNotification(webhookUrl string, msg *SlackRequestBody) error {
	cfg := types.SlackMsgEnvConfig()

	// only send slack notifications if enabled
	if cfg.Enabled == false {
		return nil
	}
	slackBody, _ := json.Marshal(msg)

	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &pooling.PoolClient{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		log.Print("unable to buffer body in slackNotification")
		return errors.New("unable to buffer body in slackNotification")
	}
	if buf.String() != "ok" {
		return errors.New("the message didn't go to slack")
	}
	return nil
}
