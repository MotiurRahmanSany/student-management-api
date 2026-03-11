package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/MotiurRahmanSany/student-management-api/internal/api/middleware"
	"github.com/MotiurRahmanSany/student-management-api/internal/response"
	"github.com/MotiurRahmanSany/student-management-api/internal/service"
)

type EnrollStudentRequest struct {
	CourseID int64 `json:"course_id"`
}

type EnrollmentHandler struct {
	enrollmentService service.EnrollmentService
}

func NewEnrollmentHandler(enrollmentService service.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{enrollmentService: enrollmentService}
}

func (h *EnrollmentHandler) EnrollStudent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || userID == "" {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var req EnrollStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if req.CourseID <= 0 {
		_ = response.Error(w, http.StatusBadRequest, "Invalid course ID", nil)
		return
	}

	enrollment, err := h.enrollmentService.EnrollStudent(r.Context(), userID, req.CourseID)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrStudentNotFound):
			_ = response.Error(w, http.StatusNotFound, "Student profile not found. Please complete your profile first.", nil)
		case errors.Is(err, service.ErrCourseNotFound):
			_ = response.Error(w, http.StatusNotFound, "The requested course does not exist.", nil)
		case errors.Is(err, service.ErrCourseAtCapacity):
			_ = response.Error(w, http.StatusConflict, "This course is already full.", nil)
		case errors.Is(err, service.ErrAlreadyEnrolled):
			_ = response.Error(w, http.StatusConflict, "You are already enrolled in this course.", nil)
		default:
			_ = response.Error(w, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	_ = response.Success(w, http.StatusCreated, "Enrolled successfully", enrollment)
}

func (h *EnrollmentHandler) UnenrollStudent(w http.ResponseWriter, r *http.Request) {
	callerUserID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || callerUserID == "" {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	callerRole, ok := r.Context().Value(middleware.RoleContextKey).(string)
	if !ok || callerRole == "" {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	studentId, err := strconv.ParseInt(r.PathValue("studentID"), 10, 64)

	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}

	courseId, err := strconv.ParseInt(r.PathValue("courseID"), 10, 64)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid course ID", nil)
		return
	}

	err = h.enrollmentService.UnenrollStudent(r.Context(), studentId, courseId, callerUserID, callerRole)

	if err != nil {
		if errors.Is(err, service.ErrCourseNotFound) {
			_ = response.Error(w, http.StatusNotFound, "Course not found", nil)
			return
		}
		if errors.Is(err, service.ErrStudentNotFound) {
			_ = response.Error(w, http.StatusNotFound, "Student not found", nil)
			return
		}

		if errors.Is(err, service.ErrNotOwner) {
			_ = response.Error(w, http.StatusForbidden, "You are not authorized to perform this action", nil)
			return
		}

		_ = response.Error(w, http.StatusInternalServerError, "An unexpected error occurred", err)
		return
	}
	
	_ = response.Success(w, http.StatusOK, "Student unenrolled successfully", nil)

}

func (h *EnrollmentHandler) ListCoursesByStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid student ID", nil)
		return
	}

	courses, err := h.enrollmentService.ListCoursesByStudent(r.Context(), id)
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, "An unexpected error occurred", err)
		return
	}

	_ = response.Success(w, http.StatusOK, "Courses retrieved successfully", courses)
}

func (h *EnrollmentHandler) ListStudentsByCourse(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid course ID", nil)
		return
	}

	students, err := h.enrollmentService.ListStudentsByCourse(r.Context(), id)
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, "An unexpected error occurred", err)
		return
	}

	_ = response.Success(w, http.StatusOK, "Students retrieved successfully", students)
}
