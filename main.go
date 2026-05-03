package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: taskctl <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  list\t\t - Show all tasks")
		fmt.Println("  add\t\t - Add new task")
		fmt.Println("  done <id>\t\t - Mark a task as completed")
	}

	connStr := "postgres://postgres:12345678@localhost:5432/taskctl"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	defer conn.Close(context.Background())

	command := os.Args[1]

	switch command {
	case "list":
		listTasks(conn)
	case "add":
		if len(os.Args) < 3 {
			log.Fatal("Please provide a task description")
		}
		addTask(conn, os.Args[2])
	case "done":
		if len(os.Args) < 3 {
			log.Fatal("Please provide a task ID")
		}
		id := 0
		fmt.Sscanf(os.Args[2], "%d", &id)
		completeTask(conn, id)
	default:
		log.Fatal("Unknown command:", command)
	}
}

func listTasks(conn *pgx.Conn) {
	rows, err := conn.Query(
		context.Background(),
		"SELECT id, description, completed FROM tasks ORDER BY id",
	)
	if err != nil {
		log.Fatal("Failed to query tasks:", err)
	}
	defer rows.Close()

	fmt.Println("Tasks:")
	for rows.Next() {
		var id int
		var description string
		var completed bool

		err = rows.Scan(&id, &description, &completed)
		if err != nil {
			log.Fatal("Failed to scan row:", err)
		}

		status := "[ ]"
		if completed {
			status = "[o]"
		}

		fmt.Printf("%d. %s %s\n", id, status, description)
	}
}

func addTask(conn *pgx.Conn, description string) {
	var id int
	err := conn.QueryRow(
		context.Background(),
		"INSERT INTO tasks (description) VALUES ($1) RETURNING id",
		description,
	).Scan(&id)

	if err != nil {
		log.Fatal("Failed to insert task:", err)
	}

	fmt.Printf("Task added with ID: %d\n", id)
}

func completeTask(conn *pgx.Conn, id int) {
	result, err := conn.Exec(
		context.Background(),
		"UPDATE tasks SET completed = true WHERE id = $1",
		id,
	)

	if err != nil {
		log.Fatal("Failed to complete task:", err)
	}
	if result.RowsAffected() == 0 {
		fmt.Printf("Task not found with ID: %d\n", id)
	}

	fmt.Printf("Task marked as completed with ID: %d", id)
}
