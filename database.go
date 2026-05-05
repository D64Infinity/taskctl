package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const (
	PENDING   = "pending"
	COMPLETED = "completed"
)

var listCommandConditions = map[string]string{
	PENDING:   "completed = false",
	COMPLETED: "completed = true",
}

func listTasks(conn *pgx.Conn, flags map[string]*bool) error {
	queryStr := "SELECT id, description, completed FROM tasks"

	for flagKey, flagValue := range flags {
		if *flagValue {
			queryStr += " WHERE " + listCommandConditions[flagKey]
			break
		}
	}
	queryStr += " ORDER BY id"

	rows, err := conn.Query(
		context.Background(),
		queryStr,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var description string
		var completed bool

		err = rows.Scan(&id, &description, &completed)
		if err != nil {
			return err
		}

		status := "[ ]"
		if completed {
			status = "[o]"
		}
		fmt.Printf("%d. %s %s\n", id, status, description)
	}
	return nil
}

func addTask(conn *pgx.Conn, description string) (int, error) {
	var id int

	err := conn.QueryRow(
		context.Background(),
		"INSERT INTO tasks (description) VALUES ($1) RETURNING id",
		description,
	).Scan(&id)

	return id, err
}

func completeTask(conn *pgx.Conn, id int) error {
	result, err := conn.Exec(
		context.Background(),
		"UPDATE tasks SET completed = true WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("Task not found with ID: %d", id)
	}
	return nil
}

func deleteTask(conn *pgx.Conn, id int) error {
	result, err := conn.Exec(
		context.Background(),
		"DELETE FROM tasks WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("Task not found with ID: %d", id)
	}
	return nil
}

func countTasks(conn *pgx.Conn) (int, int, int, error) {
	var total, active, completed int

	err := conn.QueryRow(
		context.Background(),
		`SELECT
			COUNT(id) AS total,
			COUNT(id) FILTER(WHERE NOT completed) AS active,
			COUNT(id) FILTER(WHERE completed) AS completed
		FROM tasks`,
	).Scan(&total, &active, &completed)

	return total, active, completed, err
}
