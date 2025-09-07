package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AlexFox86/scheduler/internal/api/dto"
	"github.com/AlexFox86/scheduler/internal/models"
	"github.com/AlexFox86/scheduler/internal/service/auth"
	"github.com/AlexFox86/scheduler/internal/service/tasks"
)

const dateFmt = "20060102"

// Handler provides HTTP handlers for authentication
type Handler struct {
	tasks *tasks.Service
	auth  *auth.Service
}

// NewHandler creates a new Handler
func NewHandler(tasks *tasks.Service, auth *auth.Service) *Handler {
	return &Handler{
		tasks: tasks,
		auth:  auth,
	}
}

func (h *Handler) writeData(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NextDayHandler processes the 'GET api/nextdate' request
func (h *Handler) NextDayHandler(w http.ResponseWriter, r *http.Request) {
	repeatForm := r.FormValue("repeat")
	nowForm := r.FormValue("now")
	dateForm := r.FormValue("date")

	var err error
	var now time.Time

	if nowForm == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFmt, nowForm)
		if err != nil {
			http.Error(w, "invalid 'now' parameter", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := h.tasks.NextDate(now, dateForm, repeatForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}

// AddTaskHandler processes the 'POST /api/task' request
func (h *Handler) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.writeData(w, dto.ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := h.tasks.AddTask(task)
	if err != nil {
		h.writeData(w, dto.ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	h.writeData(w, dto.ResponseID{ID: id}, http.StatusOK)
}

// GetTasksHandler processes the 'GET /api/tasks' request
func (h *Handler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := h.tasks.GetTasks(search, 50)
	if err != nil {
		h.writeData(w, err, http.StatusBadRequest)
		return
	}
	h.writeData(w, dto.TasksResp{
		Tasks: tasks,
	}, http.StatusOK)
}

// GetTaskHandler processes the 'GET /api/task' request
func (h *Handler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := h.tasks.GetTask(id)
	if err != nil {
		h.writeData(w, dto.ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	h.writeData(w, task, http.StatusOK)
}

// UpdateTaskHandler processes the 'PUT /api/task' request
func (h *Handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.writeData(w, dto.ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	err := h.tasks.Update(&task)
	if err != nil {
		h.writeData(w, dto.ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	h.writeData(w, dto.EmptyResponse{}, http.StatusOK)
}

// DoneTaskHandler processes the 'POST /api/task/done' request
func (h *Handler) DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := h.tasks.DoneTask(id)
	if err != nil {
		h.writeData(w, dto.ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	h.writeData(w, dto.EmptyResponse{}, http.StatusOK)
}

// DeleteTaskHandler processes the 'DELETE /api/task' request
func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := h.tasks.DeleteTask(id)
	if err != nil {
		h.writeData(w, dto.ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	h.writeData(w, dto.EmptyResponse{}, http.StatusOK)
}

// LoginHandler processes the 'POST /api/signin' request
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		h.writeData(w, dto.ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	token, err := h.auth.Login(login.Password)
	if err != nil {
		h.writeData(w, dto.ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	h.writeData(w, dto.ResponseToken{Token: token}, http.StatusOK)
}
