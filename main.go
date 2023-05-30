package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Complete bool   `json:"complete"`
}

var db *sql.DB

func main() {
	// Connect to the PostgreSQL database
	connStr := "user=your_username password=your_password dbname=your_database sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the tasks table if it doesn't exist
	createTable()

	// Initialize the router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	// Start the server
	log.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func createTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		complete BOOLEAN NOT NULL
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, complete FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Title, &task.Complete)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	json.NewDecoder(r.Body).Decode(&task)

	insertSQL := `INSERT INTO tasks (title, complete) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(insertSQL, task.Title, task.Complete).Scan(&task.ID)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskID := params["id"]

	var task Task
	json.NewDecoder(r.Body).Decode(&task)

	updateSQL := `UPDATE tasks SET title = $1, complete = $2 WHERE id = $3`
	_, err := db.Exec(updateSQL, task.Title, task.Complete, taskID)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskID := params["id"]

	deleteSQL := `DELETE FROM tasks WHERE id = $1`
	_, err := db.Exec(deleteSQL, taskID)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusNoContent)
}
