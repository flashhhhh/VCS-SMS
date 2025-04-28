package handler

import (
	"mail_service/internal/service"
	"net/http"
	"strconv"
)

type MailHandler interface {
	ManualSendEmail(w http.ResponseWriter, r *http.Request)
}

type mailHandler struct {
	mailService service.MailService
}

func NewMailHandler(mailService service.MailService) MailHandler {
	return &mailHandler{
		mailService: mailService,
	}
}

func (h *mailHandler) ManualSendEmail(w http.ResponseWriter, r *http.Request) {
	// Extract parameters from the request
	startTime := r.URL.Query().Get("start_time")
	endTime := r.URL.Query().Get("end_time")

	// Validate parameters
	if startTime == "" || endTime == "" {
		http.Error(w, "Missing start_time or end_time parameter", http.StatusBadRequest)
		return
	}

	// Convert to int64, assuming startTime and endTime are in Unix timestamp format
	startTimeInt, err := strconv.ParseInt(startTime, 10, 64)
	if err != nil {
		http.Error(w, "Invalid start_time format", http.StatusBadRequest)
		return
	}

	endTimeInt, err := strconv.ParseInt(endTime, 10, 64)
	if err != nil {
		http.Error(w, "Invalid end_time format", http.StatusBadRequest)
		return
	}

	// Call the mail service to send emails
	err = h.mailService.StartEmailReport(startTimeInt, endTimeInt)
	if err != nil {
		http.Error(w, "Failed to send emails: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Emails sent successfully"))
}