package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
	connStr := "postgres://postgres:12345678@localhost:5432/taskctl"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	defer conn.Close(context.Background())

	taskDescription := "Practice Golang"

	var taskId int
	err = conn.QueryRow(
		context.Background(),
		"INSERT INTO tasks (description) VALUES ($1) RETURNING id",
		taskDescription,
	).Scan(&taskId)

	if err != nil {
		log.Fatal("Failed to insert the task: ", err)
	}

	fmt.Printf("Task added with ID: %d\n", taskId)
}
