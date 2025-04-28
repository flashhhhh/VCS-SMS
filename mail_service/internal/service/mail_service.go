package service

import (
	"fmt"
	grpcclient "mail_service/internal/grpc_client"
	"time"

	"github.com/flashhhhh/pkg/env"
	"gopkg.in/gomail.v2"
)

type MailService interface {
	StartEmailReport(startTime int64, endTime int64) (error)
	PrepareEmail(to string, subject string, numServers, numOnServers, numOffServers int, meanUptimeRate float64) (error)
	SendEmail(to string, subject string, body string) error
}

type mailService struct{
	grpcClient grpcclient.ServerAdministrationServiceClient
}

func NewMailService(grpcClient grpcclient.ServerAdministrationServiceClient) MailService {
	return &mailService{
		grpcClient: grpcClient,
	}
}

func (mail *mailService) StartEmailReport(startTime int64, endTime int64) (error) {
	resp, err := mail.grpcClient.GetServerInformation(startTime, endTime)
	if err != nil {
		return err
	}

	numServers := int(resp.NumServers)
	numOnServers := int(resp.NumOnServers)
	numOffServers := int(resp.NumOffServers)
	meanUptimeRate := float64(resp.MeanUptimeRatio)

	to := env.GetEnv("SERVER_ADMINISTRATOR_EMAIL", "")
	subject := "Daily Server Status Report for " + time.Now().Format("2006-01-02")

	return mail.PrepareEmail(to, subject, numServers, numOnServers, numOffServers, meanUptimeRate)
}

func (mail *mailService) PrepareEmail(to string, subject string, numServers, numOnServers, numOffServers int, meanUptimeRate float64) (error) {
	body := fmt.Sprintf("Dear server administrator,\n\nThe server status is as follows:\n\nTotal servers: %d\nServers on: %d\nServers off: %d\nMean uptime rate: %.2f%%\n\nBest regards,\nYour Server Monitoring System", numServers, numOnServers, numOffServers, meanUptimeRate * 100)
	return mail.SendEmail(to, subject, body)
}

func (mail *mailService) SendEmail(to string, subject string, body string) error {
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