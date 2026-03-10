package router

import (
	"net/http"

	"github.com/MotiurRahmanSany/student-management-api/internal/api/handlers"
	"github.com/MotiurRahmanSany/student-management-api/internal/api/middleware"
	"github.com/MotiurRahmanSany/student-management-api/internal/auth"
)

func Setup(
	jwtManager *auth.JWTManager,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	studentHandler *handlers.StudentHandler,
	courseHandler *handlers.CourseHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Public Routes

	// health
	mux.HandleFunc("GET /health", healthHandler.Check)

	// auth
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.RefreshToken)

	// Protected Routes - require authentication
	
	authMw := middleware.AuthMiddleware(jwtManager)
	
	mux.Handle("GET /auth/me", authMw(http.HandlerFunc(authHandler.GetMe)))
	mux.Handle("POST /auth/logout", authMw(http.HandlerFunc(authHandler.Logout)))

	// student
	mux.Handle("POST /students", authMw(http.HandlerFunc(studentHandler.RegisterStudent)))
	mux.Handle("GET /students", authMw(http.HandlerFunc(studentHandler.ListAllStudents)))
	mux.Handle("GET /students/{id}", authMw(http.HandlerFunc(studentHandler.GetStudentByID))) // Expecting /students/{id}
	mux.Handle("PATCH /students/{id}", authMw(http.HandlerFunc(studentHandler.UpdateStudent))) // Expecting /students/{id}
	mux.Handle("DELETE /students/{id}", authMw(http.HandlerFunc(studentHandler.DeleteStudent))) // Expecting /students/{id} 
	
	// course
	mux.Handle("POST /courses", authMw(http.HandlerFunc(courseHandler.CreateCourse)))
	mux.Handle("GET /courses", authMw(http.HandlerFunc(courseHandler.ListAllCourses)))
	mux.Handle("GET /courses/{id}", authMw(http.HandlerFunc(courseHandler.GetCourseByID))) // Expecting /courses/{id}
	mux.Handle("PATCH /courses/{id}", authMw(http.HandlerFunc(courseHandler.UpdateCourse))) // Expecting /courses/{id}
	mux.Handle("DELETE /courses/{id}", authMw(http.HandlerFunc(courseHandler.DeleteCourse))) // Expecting /courses/{id}
	
	return mux
}
