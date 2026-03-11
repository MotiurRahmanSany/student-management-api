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
	enrollmentHandler *handlers.EnrollmentHandler,
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
	adminOnly := middleware.AdminOnly

	mux.Handle("GET /auth/me", authMw(http.HandlerFunc(authHandler.GetMe)))
	mux.Handle("POST /auth/logout", authMw(http.HandlerFunc(authHandler.Logout)))

	// student
	mux.Handle("POST /students", authMw(http.HandlerFunc(studentHandler.RegisterStudent)))
	mux.Handle("GET /students", authMw(adminOnly(http.HandlerFunc(studentHandler.ListAllStudents)))) // Only admin can list all students
	mux.Handle("GET /students/me", authMw(http.HandlerFunc(studentHandler.GetStudentByUserID)))
	mux.Handle("GET /students/{id}", authMw(adminOnly(http.HandlerFunc(studentHandler.GetStudentByID))))
	mux.Handle("PATCH /students/{id}", authMw(adminOnly(http.HandlerFunc(studentHandler.UpdateStudent))))  // Only admin can update student
	mux.Handle("DELETE /students/{id}", authMw(adminOnly(http.HandlerFunc(studentHandler.DeleteStudent)))) // Only admin can delete student

	// course
	mux.Handle("POST /courses", authMw(adminOnly(http.HandlerFunc(courseHandler.CreateCourse)))) // Only admin can create course
	mux.Handle("GET /courses", authMw(http.HandlerFunc(courseHandler.ListAllCourses)))
	mux.Handle("GET /courses/{id}", authMw(http.HandlerFunc(courseHandler.GetCourseByID)))              // Expecting /courses/{id}
	mux.Handle("PATCH /courses/{id}", authMw(adminOnly(http.HandlerFunc(courseHandler.UpdateCourse))))  // Only admin can update course
	mux.Handle("DELETE /courses/{id}", authMw(adminOnly(http.HandlerFunc(courseHandler.DeleteCourse)))) // Only admin can delete course

	// enrollment
	mux.Handle("POST /enrollments",
		authMw(adminOnly(http.HandlerFunc(enrollmentHandler.EnrollStudent))))
	mux.Handle("DELETE /enrollments/{studentID}/{courseID}",
		authMw(http.HandlerFunc(enrollmentHandler.UnenrollStudent)))

	mux.Handle("GET /students/{id}/courses",
		authMw(http.HandlerFunc(enrollmentHandler.ListCoursesByStudent)))

	mux.Handle("GET /courses/{id}/students",
		authMw(adminOnly(http.HandlerFunc(enrollmentHandler.ListStudentsByCourse))))

	return mux
}
