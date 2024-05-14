package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-gomail/gomail"
	"net/http"
	"notifications/types"
	"time"
)

// NotificationsEmalChannel Define a global channel
var NotificationsEmalChannel = make(chan *gomail.Message)

// NotificationsSlackChannel Define a global channel
var NotificationsSlackChannel = make(chan string)

// LoginHandler Handler for user login
func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var user types.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Retrieve the stored user data from the database
		var storedUser types.User
		storedUser.ID = 12
		storedUser.Username = "test"
		storedUser.Password = "test"

		// check the plain strings for the demo if they match
		if user.Password != storedUser.Password {
			fmt.Fprintf(w, "%v: %v\n", user.Password, storedUser.Password)
			http.Error(w, "Invalid email or password. handler.go hardcoded storedUser.Username, storedUser.Password", http.StatusUnauthorized)
			return
		}

		// Generate a JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, types.JWTClaims{
			UserID: storedUser.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Set token expiration time
			},
		})

		// Sign the token with the secret key
		tokenString, err := token.SignedString([]byte(types.JWTSecretKey))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the token as the response
		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenString,
		})
	}
}

// Middleware to validate JWT token and extract user ID
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Parse the token
		token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(types.JWTSecretKey), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Verify the token is valid
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract the user ID from the token claims
		claims, ok := token.Claims.(*types.JWTClaims)
		if !ok {
			http.Error(w, "Failed to extract user ID", http.StatusInternalServerError)
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

// AddMessage to validate JWT token and extract user ID
func AddMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var messagesToAdd []types.NotificationMessage
		err := json.NewDecoder(r.Body).Decode(&messagesToAdd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, message := range messagesToAdd {
			cfg := types.EmailClientEnvConfig()
			msg := gomail.NewMessage()
			msg.SetHeader("From", cfg.User)
			msg.SetHeader("To", cfg.User)
			msg.SetHeader("Subject", message.Title)
			msg.SetBody("text/plain", message.Message)

			//for demo we add to both channels
			NotificationsEmalChannel <- msg
			NotificationsSlackChannel <- message.Title + ":" + message.Message
		}
	}
}
