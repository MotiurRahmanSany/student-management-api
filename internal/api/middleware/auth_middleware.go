package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MotiurRahmanSany/student-management-api/internal/auth"
	"github.com/MotiurRahmanSany/student-management-api/internal/response"
)

type contextKey string

const (
	UserContextKey contextKey = "userID"
)

func AuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || (!strings.HasPrefix(authHeader, "Bearer ")) {
				response.Error(w, http.StatusUnauthorized, "Invalid authorization header", nil)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtManager.Verify(tokenStr)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Invalid token", nil)
				return

			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
