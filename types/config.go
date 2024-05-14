package types

func SlackMsgEnvConfig() *SlackMsgEnvVarConfig {
	return &SlackMsgEnvVarConfig{
		SlackReportWebhook: "SLACK_REPORTS_WEBHOOK",
		ENV:                "DEPLOYMENT_VAR",
		Enabled:            false,
	}
}

func EmailClientEnvConfig() *EmailClient {
	return &EmailClient{
		Host: "Host",
		Port: 0,
		User: "UN",
		Pass: "PASSWORD",
	}
}
