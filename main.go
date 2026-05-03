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
