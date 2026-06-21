package storage

import (
	"context"
	"fmt"
	"todo/internal/model"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
}

func NewStorage(dbConn *pgx.Conn) *Storage {
	return &Storage{db: dbConn}
}

func (s *Storage) GetTasks(ctx context.Context, userID int64) ([]model.ToDo, error) {
	query := "SELECT id,user_id,title,status FROM todo_list WHERE user_id = $1"

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("Storage get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []model.ToDo

	for rows.Next() {
		var t model.ToDo
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Status); err != nil {
			return nil, fmt.Errorf("Storage scan task: %w", err)
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("storage rows error: %w", err)
	}
	return tasks, nil
}

func (s *Storage) CreateTask(ctx context.Context, task model.ToDo) (int64, error) {
	query := "INSERT INTO todo_list (user_id,title,status)  VALUES ($1,$2,$3) RETURNING id"
	var lastInsert int64

	err := s.db.QueryRow(ctx, query, task.UserID, task.Title, task.Status).Scan(&lastInsert)
	if err != nil {
		return 0, fmt.Errorf("storage create task:%w", err)
	}
	return lastInsert, nil
}

func (s *Storage) UpdateTask(ctx context.Context, task model.ToDo) error {
	query := "UPDATE todo_list SET title = $1, status = $2 WHERE id = $3 AND user_id = $4"
	_, err := s.db.Exec(ctx, query, task.Title, task.Status, task.ID, task.UserID)
	if err != nil {
		return fmt.Errorf("storage update task : %w", err)
	}
	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, ID int64, userID int64) error {
	query := "DELETE FROM todo_list WHERE id = $1 AND user_id = $2"
	_, err := s.db.Exec(ctx, query, ID, userID)
	if err != nil {
		return fmt.Errorf("storage delete task : %w", err)
	}
	return nil
}

func (s *Storage) CreateUser(ctx context.Context, username string, passwordhash string) (int64, error) {
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"
	var lastInsertUser int64
	err := s.db.QueryRow(ctx, query, username, passwordhash).Scan(&lastInsertUser)
	if err != nil {
		return 0, fmt.Errorf("storage create user:%w", err)
	}
	return lastInsertUser, nil

}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	query := "SELECT id,username,password_hash FROM users WHERE username = $1"
	var u model.User

	err := s.db.QueryRow(ctx, query, username).Scan(&u.ID, &u.UserName, &u.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("storage get user by username: %w", err)
	}

	return u, nil

}
