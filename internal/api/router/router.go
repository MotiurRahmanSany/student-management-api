package router

import (
	"net/http"

	"github.com/MotiurRahmanSany/student-management-api/internal/api/handlers"
	"github.com/MotiurRahmanSany/student-management-api/internal/api/middleware"
	"github.com/MotiurRahmanSany/student-management-api/internal/auth"
)

func Setup(
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	jwtManager *auth.JWTManager,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler.Check)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.RefreshToken)

	// Protected Routes - require authentication
	authMw := middleware.AuthMiddleware(jwtManager)
	mux.Handle("GET /auth/me", authMw(http.HandlerFunc(authHandler.GetMe)))
	mux.Handle("POST /auth/logout", authMw(http.HandlerFunc(authHandler.Logout)))

	return mux
}
