package router

import (
	"net/http"

	"github.com/MotiurRahmanSany/student-management-api/internal/api/handlers"
	"github.com/MotiurRahmanSany/student-management-api/internal/api/middleware"
	"github.com/MotiurRahmanSany/student-management-api/internal/auth"
	"github.com/MotiurRahmanSany/student-management-api/internal/response"
)

func Setup(authHandler *handlers.AuthHandler, jwtManager *auth.JWTManager) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		_ = response.Success(w, http.StatusOK, "Health check successful", nil)
	})
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.RefreshToken)

	// Protected Routes - require authentication
	authMw := middleware.AuthMiddleware(jwtManager)
	mux.Handle("GET /auth/me", authMw(http.HandlerFunc(authHandler.GetMe)))
	mux.Handle("POST /auth/logout", authMw(http.HandlerFunc(authHandler.Logout)))

	return mux
}
