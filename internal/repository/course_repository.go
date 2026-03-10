package repository

import (
	"context"

	"github.com/MotiurRahmanSany/student-management-api/internal/db"
	"github.com/MotiurRahmanSany/student-management-api/internal/domain"
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, title, code string, credit, capacity int32) (domain.Course, error)
	ListCourses(ctx context.Context, limit, offset int32) ([]domain.Course, error)
	GetCourseByID(ctx context.Context, id int64) (domain.Course, error)
	UpdateCourse(ctx context.Context, id int64, title, code string, credit, capacity int32) (domain.Course, error)
	DeleteCourse(ctx context.Context, id int64) error
}

type courseRepository struct {
	q *db.Queries
}

func NewCourseRepository(q *db.Queries) CourseRepository {
	return &courseRepository{q: q}
}

func (r *courseRepository) CreateCourse(ctx context.Context, title, code string, credit, capacity int32) (domain.Course, error) {
	row, err := r.q.CreateCourse(
		ctx, db.CreateCourseParams{
			Title: title,
			Code: code,
			Credit: credit,
			Capacity: capacity,
		},
	)
	if err != nil {
		return domain.Course{}, err
	}

	course := domain.Course{
		ID:        row.ID,
		Title:     row.Title,
		Code:      row.Code,
		Credit:    row.Credit,
		Capacity:  row.Capacity,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return course, nil
}

func (r *courseRepository) ListCourses(ctx context.Context, limit, offset int32) ([]domain.Course, error) {
	rows, err := r.q.ListCourses(ctx, db.ListCoursesParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, err
	}

	courses := make([]domain.Course, len(rows))

	for i, row := range rows {
		courses[i] = domain.Course{
			ID:        row.ID,
			Title:     row.Title,
			Code:      row.Code,
			Credit:    row.Credit,
			Capacity:  row.Capacity,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
	}

	return courses, nil
}

func (r *courseRepository) GetCourseByID(ctx context.Context, id int64) (domain.Course, error) {
	row, err := r.q.GetCourseByID(ctx, id)

	if err != nil {
		return domain.Course{}, err
	}

	course := domain.Course{
		ID:        row.ID,
		Title:     row.Title,
		Code:      row.Code,
		Credit:    row.Credit,
		Capacity:  row.Capacity,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return course, nil
}

func (r *courseRepository) UpdateCourse(ctx context.Context, id int64, title, code string, credit, capacity int32) (domain.Course, error) {
	row, err := r.q.UpdateCourse(ctx, db.UpdateCourseParams{
		ID: id,
		Title:    title,
		Code:     code,
		Credit:   credit,
		Capacity: capacity,
	})

	if err != nil {
		return domain.Course{}, err
	}

	course := domain.Course{
		ID:        row.ID,
		Title:     row.Title,
		Code:      row.Code,
		Credit:    row.Credit,
		Capacity:  row.Capacity,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return course, nil
}

func (r *courseRepository) DeleteCourse(ctx context.Context, id int64) error {
	return r.q.DeleteCourse(ctx, id)

}
