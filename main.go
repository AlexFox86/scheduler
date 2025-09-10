package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AlexFox86/scheduler/internal/api"
	"github.com/AlexFox86/scheduler/internal/repository/sqlite"
	"github.com/AlexFox86/scheduler/internal/server"
	"github.com/AlexFox86/scheduler/internal/service/auth"
	"github.com/AlexFox86/scheduler/internal/service/tasks"

	_ "modernc.org/sqlite"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		fmt.Println("Env variable TODO_DBFILE not found")
		dbFile = "scheduler.db"
	}

	// Init repository
	repo := sqlite.NewSQLiteRepo()
	err := repo.Init(dbFile)
	defer repo.Close()

	if err != nil {
		panic(err)
	}

	// Create a new scheduler service
	taskService := tasks.New(repo)

	// Create a new auth service
	authService := auth.New(os.Getenv("SECRET"))

	// Init handlers and api module
	handler := api.NewHandler(taskService, authService)

	// Start server
	log.Fatal(server.Start(handler))
}
