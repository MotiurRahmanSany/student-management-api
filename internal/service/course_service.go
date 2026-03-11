package service

import (
	"context"
	"errors"

	"github.com/MotiurRahmanSany/student-management-api/internal/domain"
	"github.com/MotiurRahmanSany/student-management-api/internal/repository"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNoFieldsProvidedForCreate = errors.New("No field provided")
	ErrCourseNotFound            = errors.New("course not found")
	ErrCourseCodeAlreadyExists   = errors.New("course code already exists")
	ErrCourseTitleAlreadyExists  = errors.New("course title already exists")
)

type CourseService interface {
	CreateCourse(ctx context.Context, title, code string, credit, capacity int32) (domain.Course, error)
	ListCourses(ctx context.Context, limit, offset int32) ([]domain.Course, error)
	CountCourses(ctx context.Context) (int64, error)
	GetCourseByID(ctx context.Context, id int64) (domain.Course, error)
	UpdateCourse(ctx context.Context, id int64, title, code *string, credit, capacity *int32) (domain.Course, error)
	DeleteCourse(ctx context.Context, id int64) error
}

type courseService struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) CreateCourse(ctx context.Context, title, code string, credit, capacity int32) (domain.Course, error) {
	createdCourse, err := s.repo.CreateCourse(ctx, title, code, credit, capacity)
	if err != nil {
		return domain.Course{}, err
	}

	return createdCourse, nil
}

func (s *courseService) ListCourses(ctx context.Context, limit, offset int32) ([]domain.Course, error) {
	return s.repo.ListCourses(ctx, limit, offset)
}

func (s *courseService) CountCourses(ctx context.Context) (int64, error) {
	return s.repo.CountCourses(ctx)
}

func (s *courseService) GetCourseByID(ctx context.Context, id int64) (domain.Course, error) {
	course, err := s.repo.GetCourseByID(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Course{}, ErrCourseNotFound
		}
		return domain.Course{}, err
	}

	return course, nil
}

func (s *courseService) UpdateCourse(ctx context.Context, id int64, title, code *string, credit, capacity *int32) (domain.Course, error) {
	course, err := s.repo.GetCourseByID(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Course{}, ErrCourseNotFound
		}
		return domain.Course{}, err
	}
	if title == nil && code == nil && credit == nil && capacity == nil {
		return domain.Course{}, ErrNoFieldsProvidedForUpdate
	}
	
	if title != nil {
		course.Title = *title
	}
	if code != nil {
		course.Code = *code
	}
	if credit != nil {
		course.Credit = *credit
	}
	if capacity != nil {
		course.Capacity = *capacity
	}

	return s.repo.UpdateCourse(ctx, id, course.Title, course.Code, course.Credit, course.Capacity)
}

func (s *courseService) DeleteCourse(ctx context.Context, id int64) error {
	_, err := s.repo.GetCourseByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCourseNotFound
		}
		return err
	}

	return s.repo.DeleteCourse(ctx, id)
}
