package notification

import (
	"github.com/robfig/cron/v3"
	"log"
	"notifications/api/handler"
	"notifications/types"
)

func CronSlackReporting() {
	cronSlackReporting := cron.New()
	cronSlackReporting.AddFunc("0 */3 * * * ", sendToSlackCron)
	cronSlackReporting.Start()

	log.Print("Cron started")
	//c.Stop()  // Stop the scheduler (does not stop any jobs already running).
}

func sendToSlackCron() {
	log.Print("Send to slack called")

	messagesToSlack := make([]*types.SlackRequestBody, 0)
	for msg := range handler.NotificationsSlackChannel {
		msgToSlack := &types.SlackRequestBody{
			Text: msg,
		}
		messagesToSlack = append(messagesToSlack, msgToSlack)
		//send to slack with webhook
		//for sendMessage := range messagesToSlack {
		//			Webhook(messagesToSlack[sendMessage])
		//		}
	}
}
