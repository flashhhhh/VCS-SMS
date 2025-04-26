package service

import (
	"fmt"

	"github.com/flashhhhh/pkg/env"
	"gopkg.in/gomail.v2"
)

func PrepareEmail(to string, subject string, numServers, numOnServers, numOffServers int, meanUptimeRate float64) (string, error) {
	body := fmt.Sprintf("Dear server administrator,\n\nThe server status is as follows:\n\nTotal servers: %d\nServers on: %d\nServers off: %d\nMean uptime rate: %.2f%%\n\nBest regards,\nYour Server Monitoring System", numServers, numOnServers, numOffServers, meanUptimeRate)

	err := SendEmail(to, subject, body)
	if err != nil {
		return "", err
	}

	return body, nil
}

func SendEmail(to string, subject string, body string) error {
	senderEmail := env.GetEnv("SENDER_EMAIL", "")
	senderPassword := env.GetEnv("SENDER_PASSWORD", "")

	m := gomail.NewMessage()

	m.SetHeader("From", senderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPassword)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("Email sent successfully")
	return nil
}