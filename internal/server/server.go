package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AlexFox86/scheduler/internal/api"
)

// Start starts the server
func Start(handler *api.Handler) error {
	handler.Init()

	port := os.Getenv("TODO_PORT")
	port = "7540"
	if port == "" {
		port = ":7540"
	} else {
		port = ":" + port
	}

	fmt.Println("Server is running on port", port)

	return http.ListenAndServe(port, nil)
}
