package domain

import (
	"time"
)

type Student struct {
	ID          int64     `json:"id"`
	UserID      string    `json:"user_id"`
	FullName    string    `json:"full_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Department  string    `json:"department"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
