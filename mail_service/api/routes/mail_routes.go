package routes

import (
	"mail_service/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, mailHandler handler.MailHandler) {
	r.HandleFunc("/mail/manual_send", mailHandler.ManualSendEmail).Methods("POST")
}