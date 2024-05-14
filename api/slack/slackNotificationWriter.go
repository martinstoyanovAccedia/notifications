package notification

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"
)

type SlackWriter struct {
	w io.Writer
}

type SlackRequestBody struct {
	Text string `json:"text"`
}

func (slackWriter SlackWriter) Write(p []byte) (int, error) {
	trimLog := strings.Split(string(p), " ")
	//Remove date and time with milliseconds from the log message
	msg := string(p)[len(trimLog[0]):]
	//check if log level is Warn or above, if not, do not proceed
	logLvl := checkLogLevel(trimLog[1])
	if logLvl < 2 {
		return 0, nil
	}

	//Build formatted message with environment + time + log message
	builder := strings.Builder{}
	builder.WriteString(time.Now().Format("2006-01-02 15:04:05") + "\n")
	builder.WriteString(msg)
	msg = builder.String()

	msgToSlack := &SlackRequestBody{
		Text: msg,
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(msgToSlack)

	Webhook(msgToSlack)

	return len(p), nil
}

func checkLogLevel(partOfLog string) int {
	switch partOfLog {
	case "[INFO]":
		return 1
	case "[WARN]":
		return 2
	case "[ERROR]":
		return 3
	case "[FATAL]":
		return 4
	default:
		return 0
	}
}
