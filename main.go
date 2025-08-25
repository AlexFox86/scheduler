package main

import (
	"log"

	"github.com/AlexFox86/scheduler/internal/server"
)

func main() {
	log.Fatal(server.Start())
}
