package repository

import (
	"context"

	"github.com/MotiurRahmanSany/student-management-api/internal/db"
	"github.com/MotiurRahmanSany/student-management-api/internal/domain"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CreateUser(ctx context.Context, email, passwordHash, role string) (domain.User, error)
}

type userRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepository{q: q}
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		ID:           row.ID.String(),
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Role:         row.Role,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, email, passwordHash, role string) (domain.User, error) {
	row, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	})
	if err != nil{
		return domain.User{}, err

	}
	createdUser := domain.User{
		ID:           row.ID.String(),
		Email:        row.Email,
		Role:         row.Role,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
	return createdUser, nil
	
}
