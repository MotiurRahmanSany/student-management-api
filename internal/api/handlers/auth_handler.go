package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MotiurRahmanSany/student-management-api/internal/api/middleware"
	"github.com/MotiurRahmanSany/student-management-api/internal/response"
	"github.com/MotiurRahmanSany/student-management-api/internal/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.Email == "" || req.Password == "" {
		_ = response.Error(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	if len(req.Password) < 6 {
		_ = response.Error(w, http.StatusBadRequest, "Password must be at least 6 characters long", nil)
		return
	}

	user, err := h.service.Register(r.Context(), req.Email, req.Password)

	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyInUse) {
			_ = response.Error(w, http.StatusConflict, "Email already in use", nil)
			return
		}
		_ = response.Error(w, http.StatusInternalServerError, "Failed to register user", nil)
		return
	}

	_ = response.Success(w, http.StatusCreated, "User registered successfully", user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req RegisterReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	loginResponse, err := h.service.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	_ = response.Success(w, http.StatusOK, "User logged in successfully", loginResponse)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || userID == "" {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	user, err := h.service.GetMe(r.Context(), userID)
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, "Failed to retrieve user", nil)
		return
	}

	_ = response.Success(w, http.StatusOK, "User retrieved successfully", user)
}
