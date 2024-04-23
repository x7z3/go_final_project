package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {

	log.Println("[INFO] starting task-manager")

	// opening sqlite database
	log.Println("[INFO] Connecting to database...")
	db, err := sql.Open("sqlite", "scheduler.db")
	checkError(err, "connection")
	defer db.Close()

	// initializing new storage
	s := NewStorage(db)

	// creating new table
	s.InitDatabase()

	// creating task service
	t := NewTaskService(s)

	// intializing handlers for web-server
	http.Handle("/", http.FileServer(http.Dir("./web/")))
	http.Handle("/api/nextdate", http.HandlerFunc(NextDate))
	http.Handle("/api/task", http.HandlerFunc(t.TaskHandler))
	http.Handle("GET /api/tasks", http.HandlerFunc(t.TasksHandler))
	http.Handle("POST /api/task/done", http.HandlerFunc(t.DoneHandler))

	// starting web-server
	log.Println("[INFO] Starting server on port 7540...")
	log.Fatal(http.ListenAndServe(":7540", nil))
}
