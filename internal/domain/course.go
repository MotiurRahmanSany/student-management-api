package domain

import "time"

type Course struct {
	ID        int64            `json:"id"`
	Title     string           `json:"title"`
	Code      string           `json:"code"`
	Credit    int32            `json:"credit"`
	Capacity  int32            `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
