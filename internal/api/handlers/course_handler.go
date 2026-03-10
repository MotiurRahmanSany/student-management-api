package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/MotiurRahmanSany/student-management-api/internal/response"
	"github.com/MotiurRahmanSany/student-management-api/internal/service"
)

type CreateCourseRequest struct {
	Title    string `json:"title"`
	Code     string `json:"code"`
	Credit   int32  `json:"credit"`
	Capacity int32  `json:"capacity"`
}

type UpdateCourseRequest struct {
	Title    *string `json:"title"`
	Code     *string `json:"code"`
	Credit   *int32  `json:"credit"`
	Capacity *int32  `json:"capacity"`
}

type CourseHandler struct {
	service service.CourseService
}

func NewCourseHandler(service service.CourseService) *CourseHandler {
	return &CourseHandler{service: service}
}

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var req CreateCourseRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	defer r.Body.Close()

	if req.Capacity <= 0 || req.Credit <= 0 || req.Title == "" || req.Code == "" {
		_ = response.Error(w, http.StatusBadRequest, "Please provide valid course details", nil)
		return
	}
	
	/*
	 example of a valid request:
		{
			"title": "Course Title",
			"code": "COURSE_CODE",
			"credit": 3,
			"capacity": 30
		}
	
	 */
	if course, err := h.service.CreateCourse(r.Context(), req.Title, req.Code, req.Credit, req.Capacity); err != nil {
		_ = response.Error(w, http.StatusInternalServerError, "Failed to create course", err)
		return
	} else {
		_ = response.Success(w, http.StatusCreated, "Course created successfully", course)
	}

}

func (h *CourseHandler) ListAllCourses(w http.ResponseWriter, r *http.Request) {
	limit := int32(10) // Default limit
	offset := int32(0) // Default offset

	if v := r.URL.Query().Get("limit"); v != "" {
		if l, err := strconv.ParseInt(v, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		if o, err := strconv.ParseInt(v, 10, 32); err == nil {
			offset = int32(o)
		}
	}

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	if limit > 100 {
		limit = 100
	}

	courses, err := h.service.ListCourses(r.Context(), limit, offset)
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, "Failed to list courses", err)
		return
	}
	_ = response.Success(w, http.StatusOK, "Courses retrieved successfully", courses)

}

func (h *CourseHandler) GetCourseByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid course ID", nil)
		return
	}
	course, err := h.service.GetCourseByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCourseNotFound) {
			_ = response.Error(w, http.StatusNotFound, "Course not found", nil)
		} else {
			_ = response.Error(w, http.StatusInternalServerError, "Failed to retrieve course", err)
		}
		return
	}
	_ = response.Success(w, http.StatusOK, "Course retrieved successfully", course)
}

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid course ID", nil)
		return
	}
	
	var req UpdateCourseRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		_ = response.Error(w, http.StatusBadRequest, "Invalid course body", nil)
		return
	}
	defer r.Body.Close()
	
	course, err := h.service.UpdateCourse(r.Context(), id, req.Title, req.Code, req.Credit, req.Capacity)
	if err != nil {
		if errors.Is(err, service.ErrCourseNotFound) {
			_ = response.Error(w, http.StatusNotFound, "Course not found", nil)
		} else {
			_ = response.Error(w, http.StatusInternalServerError, "Failed to update course", err)
		}
		return
	}
	_ = response.Success(w, http.StatusOK, "Course updated successfully", course)
}

func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid course ID", nil)
		return
	}
	if err := h.service.DeleteCourse(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrCourseNotFound) {
			_ = response.Error(w, http.StatusNotFound, "Course not found", nil)
		} else {
			_ = response.Error(w, http.StatusInternalServerError, "Failed to delete course", err)
		}
		return
	}
	_ = response.Success(w, http.StatusOK, "Course deleted successfully", nil)
}
