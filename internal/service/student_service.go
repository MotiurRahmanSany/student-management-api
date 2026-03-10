package service

import (
	"context"
	"errors"
	"time"

	"github.com/MotiurRahmanSany/student-management-api/internal/domain"
	"github.com/MotiurRahmanSany/student-management-api/internal/repository"
	"github.com/jackc/pgx/v5"
)

type StudentService interface {
	RegisterStudent(ctx context.Context, userID, fullName string, dob time.Time, dept string) (domain.Student, error)
	ListAllStudents(ctx context.Context, limit, offset int32) ([]domain.Student, error)
	GetStudentByID(ctx context.Context, id int64) (domain.Student, error)
	UpdateStudent(ctx context.Context, id int64, fullName *string, dob *time.Time, dept *string, status *string) (domain.Student, error)
	DeleteStudent(ctx context.Context, id int64) error
}

var (
	ErrStudentAlreadyExists = errors.New("Student with the given user ID already exists")
	ErrStudentNotFound      = errors.New("Student not found")
	ErrNoFieldsProvided     = errors.New("No fields provided for update")
	ErrInvalidStatus        = errors.New("Invalid status value")
)

func isValidStatus(status string) bool {
	switch status {
	case "active", "suspended", "graduated":
		return true
	}
	return false
}

type studentService struct {
	studentRepo repository.StudentRepository
}

func NewStudentService(studentRepo repository.StudentRepository) StudentService {
	return &studentService{studentRepo: studentRepo}
}

func (s *studentService) RegisterStudent(ctx context.Context, userID, fullName string, dob time.Time, dept string) (domain.Student, error) {
	_, err := s.studentRepo.GetStudentByUserID(ctx, userID)
	if err == nil {
		return domain.Student{}, ErrStudentAlreadyExists
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Student{}, err
	}

	createdStudent, err := s.studentRepo.CreateStudent(
		ctx, userID, fullName, dob, dept, "active",
	)

	if err != nil {
		return domain.Student{}, err
	}

	return createdStudent, nil
}

func (s *studentService) ListAllStudents(ctx context.Context, limit, offset int32) ([]domain.Student, error) {
	return s.studentRepo.ListStudents(ctx, limit, offset)
}

func (s *studentService) GetStudentByID(ctx context.Context, id int64) (domain.Student, error) {
	student, err := s.studentRepo.GetStudentByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Student{}, ErrStudentNotFound
		}
		return domain.Student{}, err
	}
	return student, nil
}

func (s *studentService) UpdateStudent(ctx context.Context, id int64, fullName *string, dob *time.Time, dept *string, status *string) (domain.Student, error) {

	student, err := s.studentRepo.GetStudentByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Student{}, ErrStudentNotFound
		}
		return domain.Student{}, err
	}

	if fullName == nil && dob == nil && dept == nil && status == nil {
		return domain.Student{}, ErrNoFieldsProvided
	}
	if fullName != nil {
		student.FullName = *fullName
	}

	if dob != nil {
		student.DateOfBirth = *dob
	}

	if dept != nil {
		student.Department = *dept
	}

	if status != nil {
		if !isValidStatus(*status) {
			return domain.Student{}, ErrInvalidStatus
		}
		student.Status = *status
	}

	return s.studentRepo.UpdateStudent(
		ctx,
		id,
		student.FullName,
		student.DateOfBirth,
		student.Department,
		student.Status,
	)
}

func (s *studentService) DeleteStudent(ctx context.Context, id int64) error {
	_, err := s.studentRepo.GetStudentByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrStudentNotFound
		}
		return err
	}

	return s.studentRepo.DeleteStudent(ctx, id)
}
