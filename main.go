package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: taskctl <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  list\t\t\t - Show all tasks")
		fmt.Println("  count\t\t\t - Count all tasks (summary, pending, completed)")
		fmt.Println("  add\t\t\t - Add a new task")
		fmt.Println("  done <id>\t\t - Mark a task as completed")
		fmt.Println("  delete <id>\t\t - Delete a task")
		return
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres@localhost:5432/taskctl"
	}

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	defer conn.Close(context.Background())

	command := os.Args[1]

	switch command {
	case "list":
		listFlags := flag.NewFlagSet("list", flag.ExitOnError)
		flagsList := map[string]*bool{
			PENDING:   listFlags.Bool("pending", false, "show pending tasks"),
			COMPLETED: listFlags.Bool("completed", false, "show completed tasks"),
		}
		listFlags.Parse(os.Args[2:])

		fmt.Println("Tasks:")
		err = listTasks(conn, flagsList)
		if err != nil {
			log.Fatal("Failed to list tasks: ", err)
		}
	case "count":
		taskCount, taskActiveCount, taskCompleteCount, err := countTasks(conn)
		if err != nil {
			log.Fatal("Failed to count tasks: ", err)
		}
		fmt.Printf("Summary: %d | pending: %d | completed: %d\n",
			taskCount, taskActiveCount, taskCompleteCount)
	case "add":
		if len(os.Args) < 3 {
			log.Fatal("Please provide a task description")
		}

		idAdded, err := addTask(conn, os.Args[2])
		if err != nil {
			log.Fatal("Failed to insert task: ", err)
		}
		fmt.Printf("Task added with ID: %d\n", idAdded)
	case "done":
		if len(os.Args) < 3 {
			log.Fatal("Please provide a task ID")
		}

		id := 0
		n, err := fmt.Sscanf(os.Args[2], "%d", &id)
		if err != nil || n != 1 || id <= 0 {
			log.Fatal("Invalid task ID")
		}
		err = completeTask(conn, id)
		if err != nil {
			log.Fatal("Failed to complete task: ", err)
		}
		fmt.Printf("Task marked as completed with ID: %d\n", id)
	case "delete":
		if len(os.Args) < 3 {
			log.Fatal("Please provide a task ID")
		}

		id := 0
		n, err := fmt.Sscanf(os.Args[2], "%d", &id)
		if err != nil || n != 1 || id <= 0 {
			log.Fatal("Invalid task ID")
		}
		err = deleteTask(conn, id)
		if err != nil {
			log.Fatal("Failed to delete task: ", err)
		}
		fmt.Printf("Task deleted with ID: %d\n", id)
	default:
		log.Fatal("Unknown command: ", command)
	}
}
