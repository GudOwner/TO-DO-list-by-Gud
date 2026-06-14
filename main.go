package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

func getTasks(c *gin.Context) {
	query := "SELECT id, title, status FROM todo_list"
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var tasks []ToDo
	for rows.Next() {
		var t ToDo
		if err := rows.Scan(&t.ID, &t.Title, &t.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var newTodo ToDo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO todo_list(title,status) VALUES ($1,$2) RETURNING id"

	err := db.QueryRow(context.Background(), query, newTodo.Title, newTodo.Status).Scan(&newTodo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newTodo)

}
func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	var UpdateToDo ToDo
	if err := c.ShouldBindJSON(&UpdateToDo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	UpdateToDo.ID = id
	query := "UPDATE todo_list SET title = $1, status = $2 WHERE id = $3"
	_, err = db.Exec(context.Background(), query, UpdateToDo.Title, UpdateToDo.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, UpdateToDo)
}

func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	query := "DELETE FROM todo_list WHERE id = $1"
	_, err = db.Exec(context.Background(), query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func main() {
	r := gin.Default()

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

	r.GET("/task", getTasks)
	r.PUT("/task/:id", UpdateTask)
	r.POST("/task", CreateTask)
	r.DELETE("/task/:id", DeleteTask)

	r.Run(":8080")

}
