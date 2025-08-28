package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AlexFox86/scheduler/internal/api"
)

// Start starts the server
func Start() error {
	api.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = ":7540"
	}

	fmt.Println("Server is running on port", port)

	return http.ListenAndServe(port, nil)
}
