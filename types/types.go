package types

import (
	"github.com/dgrijalva/jwt-go"
)

// User struct represents the user data stored in the MySQL database
type User struct {
	ID       int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWTClaims represents the JWT claims containing the user ID
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// JWTSecretKey is the secret key used for signing JWTs
const JWTSecretKey = "some_key"

// NotificationMessage struct represents a notification message.
type NotificationMessage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SlackMsgEnvVarConfig struct {
	SlackReportWebhook string
	ENV                string
	Enabled            bool
}

type SlackRequestBody struct {
	Text string `json:"text"`
}

type EmailClient struct {
	Host string
	Port int
	User string
	Pass string
}
