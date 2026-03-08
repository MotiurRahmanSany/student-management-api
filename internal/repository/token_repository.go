package repository

import (
	"context"

	"github.com/MotiurRahmanSany/student-management-api/internal/db"
)

type TokenRepository struct {
	q *db.Queries
}

func NewTokenRepository(q *db.Queries) *TokenRepository {
	return &TokenRepository{q: q}
}

func (r *TokenRepository) CreateToken(ctx context.Context, arg db.CreateRefreshTokenParams) error {

	_, err := r.q.CreateRefreshToken(ctx, arg)
	return err
}
