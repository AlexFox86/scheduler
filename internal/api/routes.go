package api

import (
	"net/http"
)

// Init registers all handlers
func (h *Handler) Init() {
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	// http.HandleFunc("GET /api/nextdate", h.NextDayHandler)
	// http.HandleFunc("POST /api/task", h.AddTaskHandler)
	// http.HandleFunc("GET /api/tasks", h.GetTasksHandler)
	// http.HandleFunc("GET /api/task", h.GetTaskHandler)
	// http.HandleFunc("PUT /api/task", h.UpdateTaskHandler)
	// http.HandleFunc("POST /api/task/done", h.DoneTaskHandler)
	// http.HandleFunc("DELETE /api/task", h.DeleteTaskHandler)

	http.HandleFunc("GET /api/nextdate", h.NextDayHandler)
	http.HandleFunc("POST /api/task", h.Auth(h.AddTaskHandler))
	http.HandleFunc("GET /api/task", h.Auth(h.GetTaskHandler))
	http.HandleFunc("PUT /api/task", h.Auth(h.UpdateTaskHandler))
	http.HandleFunc("DELETE /api/task", h.Auth(h.DeleteTaskHandler))
	http.HandleFunc("GET /api/tasks", h.Auth(h.GetTasksHandler))
	http.HandleFunc("POST /api/task/done", h.Auth(h.DoneTaskHandler))
	http.HandleFunc("POST /api/signin", h.LoginHandler)
}
