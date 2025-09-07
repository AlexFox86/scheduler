package dto

import "github.com/AlexFox86/scheduler/internal/models"

// AddTaskRequest request to add a task
type AddTaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// LoginRequest login request
type LoginRequest struct {
	Password string `json:"password"`
}

// ResponseID response with an ID
type ResponseID struct {
	ID string `json:"id,omitempty"`
}

// TasksResp response with a list of tasks
type TasksResp struct {
	Tasks []*models.Task `json:"tasks"`
}

// ResponseToken response with a token
type ResponseToken struct {
	Token string `json:"token,omitempty"`
}

// ErrorResponse response with a error
type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

// EmptyResponse emty response
type EmptyResponse struct {
}
