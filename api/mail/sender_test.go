package mail

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmailSenderHandler(t *testing.T) {
	// Create a mock request with dummy data
	req, err := http.NewRequest("GET", "/email", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a mock handler for the EmailSender
	handler := EmailSender()

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
