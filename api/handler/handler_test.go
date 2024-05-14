package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"notifications/api/handler"
	"notifications/types"
	"testing"
)

var token string

func TestLoginHandler(t *testing.T) {
	// Create a request body
	user := map[string]string{"username": "test", "password": "test"}
	userJSON, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.LoginHandler())

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	require.NoError(t, err)
	require.Contains(t, rr.Body.String(), "token")
}

func TestAddMessageHandler(t *testing.T) {
	// Create a sample request body
	notifications := []types.NotificationMessage{
		{Title: "Notification 1", Message: "This is message 1"},
		{Title: "Notification 2", Message: "This is message 2"},
		// Add more messages as needed
	}
	notificationsJSON, _ := json.Marshal(notifications)

	// Create a request with the token and JSON payload
	req, err := http.NewRequest("POST", "/messages", bytes.NewBuffer(notificationsJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", token)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a mock AuthMiddleware handler
	handler := http.HandlerFunc(handler.AuthMiddleware(handler.AddMessage()))

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

}
