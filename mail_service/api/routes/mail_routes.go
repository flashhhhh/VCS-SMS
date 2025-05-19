package routes

import (
	"mail_service/internal/handler"
	"mail_service/api/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, mailHandler handler.MailHandler) {
	r.Handle("/manual_send", middleware.AdminMiddleware(http.HandlerFunc(mailHandler.ManualSendEmail))).Methods("POST")
}