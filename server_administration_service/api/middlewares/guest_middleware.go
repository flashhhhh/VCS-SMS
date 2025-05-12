package middlewares

import (
	"net/http"

	"github.com/flashhhhh/pkg/jwt"
	"github.com/flashhhhh/pkg/logging"
)

func GuestMiddleware(next http.Handler) http.Handler {
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

		// Check if the user is at least a guest
		if data["role"] != "admin" && data["role"] != "user" && data["role"] != "guest" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		
		logging.LogMessage("server_administration_service", "User is a guest", "INFO")
		next.ServeHTTP(w, r)
	})
}