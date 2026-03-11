package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/MotiurRahmanSany/student-management-api/internal/config"
	"github.com/MotiurRahmanSany/student-management-api/internal/database"
	"github.com/MotiurRahmanSany/student-management-api/internal/db"
	"github.com/MotiurRahmanSany/student-management-api/internal/repository"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	adminMail := "admin@example.com"
	adminPassword := "admin123"

	cfg := config.GetConfig()

	pool, err := database.NewConnection(cfg.Db)

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	queries := db.New(pool)

	userRepo := repository.NewUserRepository(queries)

	ctx := context.Background()

	_, err = userRepo.GetUserByEmail(ctx, adminMail)
	if err == nil {
		fmt.Println("Admin already exists, skipping...")
		return
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		fmt.Println("Failed to get admin user:", err)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Failed to generate password hash:", err)
		return
	}
	user, err := userRepo.CreateUser(ctx, adminMail, string(hash), "admin")

	if err != nil {
		fmt.Println("Failed to create admin user:", err)
		return
	}

	fmt.Printf("Admin Created, email: %s (id: %s)\n", user.Email, user.ID)

}
