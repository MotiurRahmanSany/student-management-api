package service

import (
	"context"
	"errors"

	"github.com/MotiurRahmanSany/student-management-api/internal/domain"
	"github.com/MotiurRahmanSany/student-management-api/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrAlreadyEnrolled    = errors.New("student is already enrolled in this course")
	ErrCourseAtCapacity   = errors.New("course has reached maximum capacity")
	ErrEnrollmentNotFound = errors.New("enrollment not found")
	ErrNotOwner           = errors.New("you can only manage your own enrollments")
)

type EnrollmentService interface {
	EnrollStudent(ctx context.Context, userID string, courseID int64) (domain.Enrollment, error)
	UnenrollStudent(ctx context.Context, studentID, courseID int64, callerUserID, callerRole string) error
	ListCoursesByStudent(ctx context.Context, studentID int64) ([]domain.Course, error)
	ListStudentsByCourse(ctx context.Context, courseID int64) ([]domain.Student, error)
}

type enrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
	courseRepo     repository.CourseRepository
	studentRepo    repository.StudentRepository
}

func NewEnrollmentService(enrollmentRepo repository.EnrollmentRepository, courseRepo repository.CourseRepository, studentRepo repository.StudentRepository) EnrollmentService {
	return &enrollmentService{enrollmentRepo: enrollmentRepo, courseRepo: courseRepo, studentRepo: studentRepo}
}

func (s *enrollmentService) EnrollStudent(ctx context.Context, userID string, courseID int64) (domain.Enrollment, error) {
	student, err := s.studentRepo.GetStudentByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Enrollment{}, ErrStudentNotFound
		}
		return domain.Enrollment{}, err
	}

	course, err := s.courseRepo.GetCourseByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Enrollment{}, ErrCourseNotFound
		}
		return domain.Enrollment{}, err
	}

	count, err := s.enrollmentRepo.CountEnrollmentsByCourse(ctx, courseID)
	if err != nil {
		return domain.Enrollment{}, err
	}

	if count >= int64(course.Capacity) {
		return domain.Enrollment{}, ErrCourseAtCapacity
	}

	enrollment, err := s.enrollmentRepo.EnrollStudent(ctx, student.ID, courseID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.Enrollment{}, ErrAlreadyEnrolled
			}
		}
		return domain.Enrollment{}, err
	}

	return enrollment, nil

}

func (s *enrollmentService) ListCoursesByStudent(ctx context.Context, studentID int64) ([]domain.Course, error) {
	return s.enrollmentRepo.ListCoursesByStudent(ctx, studentID)
}

func (s *enrollmentService) ListStudentsByCourse(ctx context.Context, courseID int64) ([]domain.Student, error) {
	return s.enrollmentRepo.ListStudentsByCourse(ctx, courseID)
}

func (s *enrollmentService) UnenrollStudent(ctx context.Context, studentID, courseID int64, callerUserID, callerRole string) error {
	student, err := s.studentRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrStudentNotFound
		}
		return err
	}
	_, err = s.courseRepo.GetCourseByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCourseNotFound
		}
		return err
	}
	

	if callerRole != "admin" && student.UserID != callerUserID {
		return ErrNotOwner
	}

	return s.enrollmentRepo.UnenrollStudent(ctx, studentID, courseID)
}
