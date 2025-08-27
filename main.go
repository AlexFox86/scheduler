package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AlexFox86/scheduler/internal/repository/sqlite"
	"github.com/AlexFox86/scheduler/internal/server"

	_ "modernc.org/sqlite"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		fmt.Println("Env variable TODO_DBFILE not found")
		dbFile = "scheduler.db"
	}

	repo := sqlite.NewSQLiteRepo()

	err := repo.Init(dbFile)
	defer repo.Close()

	if err != nil {
		panic(err)
	}

	log.Fatal(server.Start())
}
