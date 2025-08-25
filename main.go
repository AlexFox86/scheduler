package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	//http.HandleFunc("GET /validate", handler.Validate)
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.Handle("/js/scripts.min.js", http.FileServer(http.Dir(webDir)))
	http.Handle("/css/style.css", http.FileServer(http.Dir(webDir)))
	http.Handle("/favicon.ico", http.FileServer(http.Dir(webDir)))

	port := os.Getenv("AUTH_PORT")
	port = ":7540"
	fmt.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
