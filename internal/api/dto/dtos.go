package dto

import "github.com/AlexFox86/scheduler/internal/models"

// AddTaskRequest request to add a task
type AddTaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Response response with an ID or Error
type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// TasksResp response with a list of tasks
type TasksResp struct {
	Tasks []*models.Task `json:"tasks"`
}
