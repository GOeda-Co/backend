package model

type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse is a success message response format
type MessageResponse struct {
	Message string `json:"message"`
}

// RegisterResponse example
type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// LoginResponse example
type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

// AdminCheckResponse for admin status check
type AdminCheckResponse struct {
	IsAdmin bool `json:"is_admin"`
}