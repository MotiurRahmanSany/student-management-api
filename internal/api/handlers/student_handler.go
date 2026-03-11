package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MotiurRahmanSany/student-management-api/internal/api/middleware"
	"github.com/MotiurRahmanSany/student-management-api/internal/response"
	"github.com/MotiurRahmanSany/student-management-api/internal/service"
)

type RegisterStudentReq struct {
	FullName    string `json:"full_name"`
	DateOfBirth string `json:"date_of_birth"` // Expecting ISO format: "YYYY-MM-DD"
	Department  string `json:"department"`
}

type UpdateStudentReq struct {
	FullName    *string `json:"full_name"`
	DateOfBirth *string `json:"date_of_birth"`
	Department  *string `json:"department"`
	Status      *string `json:"status"`
}

type StudentHandler struct {
	service service.StudentService
}

func NewStudentHandler(s service.StudentService) *StudentHandler {
	return &StudentHandler{
		service: s,
	}
}

func (h *StudentHandler) RegisterStudent(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || userID == "" {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req RegisterStudentReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if req.FullName == "" || req.DateOfBirth == "" || req.Department == "" {
		_ = response.Error(w, http.StatusBadRequest, "Full name, date of birth, and department are required", nil)
		return
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid date of birth format. Expected YYYY-MM-DD", nil)
		return
	}

	/*
		example:
		{
			"full_name": "John Doe",
			"date_of_birth": "2000-01-01",
			"department": "Computer Science"
		}
	*/

	student, err := h.service.RegisterStudent(r.Context(), userID, req.FullName, dob, req.Department)
	if err != nil {
		if errors.Is(err, service.ErrStudentAlreadyExists) {
			_ = response.Error(w, http.StatusConflict, service.ErrStudentAlreadyExists.Error(), nil)
		} else {
			_ = response.Error(w, http.StatusInternalServerError, "Failed to register student", err)
		}
		return
	}

	_ = response.Success(w, http.StatusCreated, "Student registered successfully", student)
}

func (h *StudentHandler) ListAllStudents(w http.ResponseWriter, r *http.Request) {
	limit := int32(10) // Default limit
	page := int32(0)   // Default offset

	if v := r.URL.Query().Get("limit"); v != "" {
		if l, err := strconv.ParseInt(v, 10, 32); err == nil {
			limit = int32(l)
		}
	}
	if v := r.URL.Query().Get("page"); v != "" {
		if p, err := strconv.ParseInt(v, 10, 32); err == nil {
			page = int32(p)
		}
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit // convert page → offset for SQL

	students, err := h.service.ListAllStudents(r.Context(), limit, offset)
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, "Failed to list students", err)
		return
	}

	total, err := h.service.CountStudents(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, "Failed to count students", err)
		return
	}
	// _ = response.Success(w, http.StatusOK, "Students retrieved successfully", students)
	_ = response.Success(w, http.StatusOK, "Students retrieved successfully",
		response.NewPaginatedData(students, page, limit, total))
}

func (h *StudentHandler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}
	student, err := h.service.GetStudentByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrStudentNotFound) {
			_ = response.Error(w, http.StatusNotFound, "Student not found", nil)
		} else {
			_ = response.Error(w, http.StatusInternalServerError, "Failed to retrieve student", err)
		}
		return
	}

	_ = response.Success(w, http.StatusOK, "Student retrieved successfully", student)
}

func (h *StudentHandler) GetStudentByUserID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || userID == "" {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	student, err := h.service.GetStudentByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrStudentProfileNotFound) {
			_ = response.Error(w, http.StatusNotFound, "Student profile not found", nil)
		} else {
			_ = response.Error(w, http.StatusInternalServerError, "Failed to retrieve student", err)
		}
		return
	}
	_ = response.Success(w, http.StatusOK, "Student retrieved successfully", student)
}

func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}

	var req UpdateStudentReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	var dob *time.Time

	if req.DateOfBirth != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			_ = response.Error(w, http.StatusBadRequest, "Invalid date format", nil)
			return
		}
		dob = &parsed
	}

	student, err := h.service.UpdateStudent(
		r.Context(),
		id,
		req.FullName,
		dob,
		req.Department,
		req.Status,
	)

	if err != nil {
		if errors.Is(err, service.ErrStudentNotFound) {
			_ = response.Error(w, http.StatusNotFound, "Student not found", nil)
			return
		}
		if errors.Is(err, service.ErrNoFieldsProvidedForUpdate) {
			_ = response.Error(w, http.StatusBadRequest, "No fields provided for update", nil)
			return
		}

		if errors.Is(err, service.ErrInvalidStatus) {
			_ = response.Error(w, http.StatusBadRequest, "Invalid status value", nil)
			return
		}

		_ = response.Error(w, http.StatusInternalServerError, "Failed to update student", err)
		return
	}

	_ = response.Success(w, http.StatusOK, "Student updated successfully", student)
}

func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}

	if err := h.service.DeleteStudent(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrStudentNotFound) {
			_ = response.Error(w, http.StatusNotFound, "Student not found", nil)
		} else {
			_ = response.Error(w, http.StatusInternalServerError, "Failed to delete student", err)
		}
		return
	}

	_ = response.Success(w, http.StatusOK, "Student deleted successfully", nil)
}
