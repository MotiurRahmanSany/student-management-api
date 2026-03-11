package repository

import (
	"context"
	"time"

	"github.com/MotiurRahmanSany/student-management-api/internal/db"
	"github.com/MotiurRahmanSany/student-management-api/internal/domain"
)

type EnrollmentRepository interface {
	EnrollStudent(ctx context.Context, studentID, courseID int64) (domain.Enrollment, error)
	ListCoursesByStudent(ctx context.Context, studentID int64) ([]domain.Course, error)
	ListStudentsByCourse(ctx context.Context, courseID int64) ([]domain.Student, error)
	CountEnrollmentsByCourse(ctx context.Context, courseID int64) (int64, error)
	UnenrollStudent(ctx context.Context, studentID, courseID int64) error
}

type enrollmentRepository struct {
	q *db.Queries
}

func NewEnrollmentRepository(q *db.Queries) EnrollmentRepository {
	return &enrollmentRepository{q: q}
}

func (r *enrollmentRepository) EnrollStudent(ctx context.Context, studentID, courseID int64) (domain.Enrollment, error) {
	row, err := r.q.EnrollStudent(ctx, db.EnrollStudentParams{
		StudentID: studentID,
		CourseID:  courseID,
	})

	if err != nil {
		return domain.Enrollment{}, err
	}

	erollment := domain.Enrollment{
		ID:         row.ID,
		StudentID:  row.StudentID,
		CourseID:   row.CourseID,
		Grade:      nil,
		EnrolledAt: row.EnrolledAt.Time,
	}

	return erollment, nil
}

func (r *enrollmentRepository) ListCoursesByStudent(ctx context.Context, studentID int64) ([]domain.Course, error) {
	rows, err := r.q.ListCoursesByStudent(ctx, studentID)

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

func (r *enrollmentRepository) ListStudentsByCourse(ctx context.Context, courseID int64) ([]domain.Student, error) {
	rows, err := r.q.ListStudentsByCourse(ctx, courseID)

	if err != nil {
		return nil, err
	}

	students := make([]domain.Student, len(rows))

	for i, row := range rows {
		students[i] = domain.Student{
			ID:          row.ID,
			UserID:      "",
			FullName:    row.FullName,
			DateOfBirth: time.Time{},
			Department:  row.Department,
			Status:      row.Status,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   time.Time{},
		}
	}

	return students, nil
}

func (r *enrollmentRepository) CountEnrollmentsByCourse(ctx context.Context, courseID int64) (int64, error) {
	count, err := r.q.CountEnrollmentsByCourse(ctx, courseID)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *enrollmentRepository) UnenrollStudent(ctx context.Context, studentID, courseID int64) error {
	return r.q.UnenrollStudent(ctx, db.UnenrollStudentParams{
		StudentID: studentID,
		CourseID:  courseID,
	})

}
