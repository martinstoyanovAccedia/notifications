package mail

import (
	"crypto/tls"
	"github.com/go-gomail/gomail"
	"log"
	"net/http"
	"notifications/api/handler"
	"notifications/types"
	"time"
)

// EmailSender to validate JWT token and extract user ID
func EmailSender() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go func() {
			cfg := types.EmailClientEnvConfig()
			d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Pass)
			d.TLSConfig = &tls.Config{
				ServerName:         "smtp.abv.bg",
				InsecureSkipVerify: true,
			}
			var s gomail.SendCloser
			var err error
			open := false
			for {
				select {
				case m, ok := <-handler.NotificationsEmalChannel:
					if !ok {
						return
					}
					if !open {
						// Dial with the Dialer
						if s, err = d.Dial(); err != nil {
							log.Printf("Error dialing SMTP server: %v", err)
							continue // Skip sending the message
						}
						open = true
					}
					if err := gomail.Send(s, m); err != nil {
						log.Print(err)
						//TODO: store unsuccessful messages somewhere
						//Scheduler to retry
					}
				// Close the connection to the SMTP server if no email was sent in
				// the last 5 minutes.
				case <-time.After(300 * time.Second):
					if open {
						if err := s.Close(); err != nil {
							panic(err)
						}
						open = false
					}
				}
			}
		}()
	}
}
