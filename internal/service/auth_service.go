package service

import (
	"context"
	"errors"

	"github.com/MotiurRahmanSany/student-management-api/internal/auth"
	"github.com/MotiurRahmanSany/student-management-api/internal/domain"
	"github.com/MotiurRahmanSany/student-management-api/internal/repository"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyInUse = errors.New("email already in use")

type LoginResponse struct {
	User        domain.User `json:"user"`
	AccessToken string      `json:"access_token"`
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (LoginResponse, error)
	GetMe(ctx context.Context, userID string) (domain.User, error)
	Logout(ctx context.Context ) error
}

type authService struct {
	userRepo repository.UserRepository
	jwt      *auth.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwt *auth.JWTManager) AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (s *authService) Register(ctx context.Context, email, password string) (domain.User, error) {
	_, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return domain.User{}, ErrEmailAlreadyInUse
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	return s.userRepo.CreateUser(ctx, email, string(hash), "student")
}

func (s *authService) Login(ctx context.Context, email, password string) (LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return LoginResponse{}, err
	}

	token, err := s.jwt.Generate(user.ID)

	if err != nil {
		return LoginResponse{}, err
	}

	res := LoginResponse{
		User:        user,
		AccessToken: token,
	}

	return res, nil

}

func (s *authService) GetMe(ctx context.Context, userID string) (domain.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func (s *authService) Logout(ctx context.Context) error{
	return nil
}