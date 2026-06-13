package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type ToDo struct {
	ID     int64  `json:"id"`
	Title  string `json:"task"`
	Status bool   `json:"status"`
}

func InitDBURL() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error loading .env %w", err)
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	if user == "" || dbname == "" {
		return "", fmt.Errorf("Not all required environment variables are set in .env")
	}

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	return dbURL, nil

}

var db *pgx.Conn

func HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTasks(w, r)
	case http.MethodPost:
		CreateTask(w, r)
	case http.MethodPut:
		UpdateTask(w, r)
	case http.MethodDelete:
		DeleteTask(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, title, status FROM todo_list"
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []ToDo
	for rows.Next() {
		var t ToDo
		if err := rows.Scan(&t.ID, &t.Title, &t.Status); err != nil {
			http.Error(w, "Ошибка чтения строки", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "Ошибка базы данных после чтения", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTodo ToDo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO todo_list(title,status) VALUES ($1,$2) RETURNING id"

	err := db.QueryRow(context.Background(), query, newTodo.Title, newTodo.Status).Scan(&newTodo.ID)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)

}
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var UpdateToDo ToDo
	if err := json.NewDecoder(r.Body).Decode(&UpdateToDo); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	query := "UPDATE todo_list SET title = $1, status = $2 WHERE id = $3"
	_, err := db.Exec(context.Background(), query, UpdateToDo.Title, UpdateToDo.Status, UpdateToDo.ID)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UpdateToDo)

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	var delTodo ToDo
	if err := json.NewDecoder(r.Body).Decode(&delTodo); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM todo_list WHERE id = $1"
	_, err := db.Exec(context.Background(), query, delTodo.ID)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dbURL, err := InitDBURL()
	if err != nil {
		log.Fatalf("Error with init cfg %v", err)
	}
	db, err = pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to db %v", err)
	}
	defer db.Close(context.Background())

	err = db.Ping(ctx)
	if err != nil {
		log.Fatalf("Connect with db is lost %v", err)
	}
	fmt.Println("Successful! db is connecting")

	http.HandleFunc("/task", HandleTasks)
	http.ListenAndServe(":8080", nil)

}
