package api

import (
	"net/http"
)

// Init registers all handlers
func (h *Handler) Init() {
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("GET /api/nextdate", h.NextDayHandler)
	http.HandleFunc("POST /api/task", h.AddTaskHandler)
	http.HandleFunc("GET /api/tasks", h.GetTasksHandler)
	http.HandleFunc("GET /api/task", h.GetTaskHandler)
	http.HandleFunc("PUT /api/task", h.UpdateTaskHandler)
}
