package main

import (
	"log"
	"net/http"
	"notifications/api/handler"
	"notifications/api/mail"
	notification "notifications/api/slack"
)

func main() {

	//cron for slack notification
	log.Println("Cron call")
	notification.CronSlackReporting()

	http.HandleFunc("/login", handler.LoginHandler())

	//auth once to email client, open
	http.HandleFunc("/email", handler.AuthMiddleware(mail.EmailSender()))

	//adds messages to both email and slack
	http.HandleFunc("/messages", handler.AuthMiddleware(handler.AddMessage()))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
