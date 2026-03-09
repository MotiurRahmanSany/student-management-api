package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MotiurRahmanSany/student-management-api/internal/api/handlers"
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

	userRepo := repository.NewUserRepository(queries)
	tokenRepo := repository.NewTokenRepository(queries)
	jwtManager := auth.NewJWTManager(config.JwtSecretKey, time.Minute*10)

	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager)

	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)

	mux := router.Setup(healthHandler, authHandler, jwtManager)

	fmt.Printf("Server is running on port %d\n", config.HttpPort)
	fmt.Printf("Base Url is: http:localhost:%d", config.HttpPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.HttpPort), mux); err != nil {
		fmt.Printf("Error starting server: %v\n", err.Error())
	}

}
