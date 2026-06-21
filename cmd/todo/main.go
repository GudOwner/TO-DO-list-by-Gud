package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"todo/internal/handler"
	"todo/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

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
		return "", fmt.Errorf("not all required environment variables are set in .env")
	}

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	return dbURL, nil
}

func main() {
	r := gin.Default()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbURL, err := InitDBURL()
	if err != nil {
		log.Fatalf("Error with init cfg %v", err)
	}

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to db %v", err)
	}
	defer conn.Close(context.Background())

	store := storage.NewStorage(conn)

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("Connect with db is lost %v", err)
	}
	fmt.Println("Successful! db is connecting")

	newHandler := handler.NewHandler(store)

	r.LoadHTMLGlob("frontend/*.html")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.html", nil)
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", newHandler.Register)
		auth.POST("/login", newHandler.Login)
	}

	api := r.Group("/api")
	api.Use(newHandler.AuthMiddleware())
	{
		api.GET("/tasks", newHandler.GetTask)
		api.POST("/tasks", newHandler.CreateTask)
		api.PUT("/tasks/:id", newHandler.UpdateTask)
		api.DELETE("/tasks/:id", newHandler.DeleteTask)
	}

	r.Run(":8080")
}
