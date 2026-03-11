package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MotiurRahmanSany/student-management-api/internal/api/handlers"
	"github.com/MotiurRahmanSany/student-management-api/internal/api/middleware"
	"github.com/MotiurRahmanSany/student-management-api/internal/api/router"
	"github.com/MotiurRahmanSany/student-management-api/internal/auth"
	"github.com/MotiurRahmanSany/student-management-api/internal/config"
	"github.com/MotiurRahmanSany/student-management-api/internal/database"
	"github.com/MotiurRahmanSany/student-management-api/internal/db"
	"github.com/MotiurRahmanSany/student-management-api/internal/repository"
	"github.com/MotiurRahmanSany/student-management-api/internal/service"
)

func serve(config *config.Config) {
	// connecting to database
	pool, err := database.NewConnection(config.Db)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err.Error())
		return
	}
	defer pool.Close()
	queries := db.New(pool)

	jwtManager := auth.NewJWTManager(config.JwtSecretKey, time.Minute*10)
	tokenRepo := repository.NewTokenRepository(queries)

	userRepo := repository.NewUserRepository(queries)
	studentRepo := repository.NewStudentRepository(queries)
	courseRepo := repository.NewCourseRepository(queries)
	enrollmentRepo := repository.NewEnrollmentRepository(queries)

	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager)
	studentService := service.NewStudentService(studentRepo)
	courseService := service.NewCourseService(courseRepo)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, courseRepo, studentRepo)

	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	studentHandler := handlers.NewStudentHandler(studentService)
	courseHandler := handlers.NewCourseHandler(courseService)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentService)

	mux := router.Setup(
		jwtManager,
		healthHandler,
		authHandler,
		studentHandler,
		courseHandler,
		enrollmentHandler,
	)

	loggedMux := middleware.Logger(mux)

	fmt.Printf("Server is running on port %d\n", config.HttpPort)
	fmt.Printf("Base Url is: http:localhost:%d\n", config.HttpPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.HttpPort), loggedMux); err != nil {
		fmt.Printf("Error starting server: %v\n", err.Error())
	}

}
