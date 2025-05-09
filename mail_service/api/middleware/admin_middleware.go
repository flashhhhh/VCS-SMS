package middleware

import (
	"net/http"

	"github.com/flashhhhh/pkg/jwt"
	"github.com/flashhhhh/pkg/logging"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from request header bearer
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token after "Bearer "
		token := authHeader[7:]

		data, err := jwt.ValidateToken(token)
		if err != nil {
			http.Error(w, "Token is invalid", http.StatusUnauthorized)
			return
		}

		// Check if the user is an admin
		if data["role"] != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		logging.LogMessage("mail_service", "Admin authenticated successfully", "INFO")
		
		next.ServeHTTP(w, r)
	})
}