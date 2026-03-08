package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    any         `json:"data,omitempty"`
	Error   any         `json:"error,omitempty"`
}

func jsonResponse(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	
	return encoder.Encode(payload)
}

func Error(w http.ResponseWriter, status int, message string, err any) error {
	response := APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
	return jsonResponse(w, status, response)
}

func Success(w http.ResponseWriter, status int, message string, data any) error {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	return jsonResponse(w, status, response)
}
