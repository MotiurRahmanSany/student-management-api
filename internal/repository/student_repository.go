package repository

import (
	"context"
	"time"

	"github.com/MotiurRahmanSany/student-management-api/internal/db"
	"github.com/MotiurRahmanSany/student-management-api/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type StudentRepository interface {
	CreateStudent(ctx context.Context, userID, fullName string, dob time.Time, dept, status string) (domain.Student, error)
	GetStudentByID(ctx context.Context, id int64) (domain.Student, error)
	GetStudentByUserID(ctx context.Context, userID string) (domain.Student, error)
	ListStudents(ctx context.Context, limit, offset int32) ([]domain.Student, error)
	CountStudents(ctx context.Context) (int64, error)
	UpdateStudent(ctx context.Context, id int64, fullName string, dob time.Time, dept, status string) (domain.Student, error)
	DeleteStudent(ctx context.Context, id int64) error
}

type studentRepository struct {
	q *db.Queries
}

func NewStudentRepository(q *db.Queries) StudentRepository {
	return &studentRepository{q: q}
}

func (r *studentRepository) CreateStudent(ctx context.Context, userID, fullName string, dob time.Time, dept, status string) (domain.Student, error) {
	var pgID pgtype.UUID
	var pgDob pgtype.Date

	err := pgID.Scan(userID)
	if err != nil {
		return domain.Student{}, err
	}

	err = pgDob.Scan(dob)
	if err != nil {
		return domain.Student{}, err
	}

	row, err := r.q.CreateStudent(ctx, db.CreateStudentParams{
		UserID:      pgID,
		FullName:    fullName,
		DateOfBirth: pgDob,
		Department:  dept,
		Status:      status,
	})

	if err != nil {
		return domain.Student{}, err
	}

	createdStudent := domain.Student{
		ID:          row.ID,
		UserID:      row.UserID.String(),
		FullName:    row.FullName,
		DateOfBirth: row.DateOfBirth.Time,
		Department:  row.Department,
		Status:      row.Status,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	return createdStudent, nil
}

func (r *studentRepository) GetStudentByID(ctx context.Context, id int64) (domain.Student, error) {
	row, err := r.q.GetStudentByID(ctx, id)
	if err != nil {
		return domain.Student{}, err
	}

	student := domain.Student{
		ID:          row.ID,
		UserID:      row.UserID.String(),
		FullName:    row.FullName,
		DateOfBirth: row.DateOfBirth.Time,
		Department:  row.Department,
		Status:      row.Status,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	return student, nil
}

func (r *studentRepository) GetStudentByUserID(ctx context.Context, userID string) (domain.Student, error) {
	var pgID pgtype.UUID

	err := pgID.Scan(userID)

	if err != nil {
		return domain.Student{}, err
	}

	row, err := r.q.GetStudentByUserID(ctx, pgID)
	if err != nil {
		return domain.Student{}, err
	}

	return domain.Student{
		ID:          row.ID,
		UserID:      row.UserID.String(),
		FullName:    row.FullName,
		DateOfBirth: row.DateOfBirth.Time,
		Department:  row.Department,
		Status:      row.Status,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}, nil
}

func (r *studentRepository) ListStudents(ctx context.Context, limit, offset int32) ([]domain.Student, error) {
	rows, err := r.q.ListStudents(ctx, db.ListStudentsParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, err
	}
	students := make([]domain.Student, len(rows))

	for i, row := range rows {
		students[i] = domain.Student{
			ID:          row.ID,
			UserID:      row.UserID.String(),
			FullName:    row.FullName,
			DateOfBirth: row.DateOfBirth.Time,
			Department:  row.Department,
			Status:      row.Status,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		}
	}

	return students, nil
}

func (r *studentRepository) CountStudents(ctx context.Context) (int64, error) {
	return r.q.CountStudents(ctx)
}

func (r *studentRepository) UpdateStudent(ctx context.Context, id int64, fullName string, dob time.Time, dept, status string) (domain.Student, error) {
	var pgDob pgtype.Date

	err := pgDob.Scan(dob)

	if err != nil {
		return domain.Student{}, err
	}

	row, err := r.q.UpdateStudent(ctx, db.UpdateStudentParams{
		ID:          id,
		FullName:    fullName,
		DateOfBirth: pgDob,
		Department:  dept,
		Status:      status,
	})

	if err != nil {
		return domain.Student{}, err
	}

	updatedStudent := domain.Student{
		ID:          row.ID,
		UserID:      row.UserID.String(),
		FullName:    row.FullName,
		DateOfBirth: row.DateOfBirth.Time,
		Department:  row.Department,
		Status:      row.Status,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	return updatedStudent, nil

}

func (r *studentRepository) DeleteStudent(ctx context.Context, id int64) error {
	return r.q.DeleteStudent(ctx, id)
}
