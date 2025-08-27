package server

import (
	"fmt"
	"net/http"
	"os"
)

// Start starts the server
func Start() error {
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = ":7540"
	}

	fmt.Println("Server is running on port", port)

	return http.ListenAndServe(port, nil)
}
