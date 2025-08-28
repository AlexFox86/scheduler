package api

import (
	"net/http"
)

// Init registers all handlers
func Init() {
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("GET /api/nextdate", NextDayHandler)
}
