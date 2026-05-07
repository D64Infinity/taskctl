package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const (
	PENDING   = "pending"
	COMPLETED = "completed"
	HIGH      = "high"
	MEDIUM    = "medium"
	LOW       = "low"
)

var listCommandConditions = map[string]string{
	PENDING:   "completed = false",
	COMPLETED: "completed = true",
}

func retrievePriorityId(conn *pgx.Conn, priorityStr string) (int, error) {
	rows, err := conn.Query(
		context.Background(),
		"SELECT id FROM task_priority WHERE priority = $1",
		priorityStr,
	)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var id int
	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return -1, err
		}
		return id, nil
	} else {
		return -1, fmt.Errorf("Priority was not found: %s", priorityStr)
	}
}

func listTasks(conn *pgx.Conn, flags map[string]*bool) error {
	queryStr := `SELECT 
					tasks.id,
					tasks.description,
					tasks.completed,
					task_priority.priority
				 FROM tasks INNER JOIN task_priority ON task_priority.id = tasks.priority_id`

	for flagKey, flagValue := range flags {
		if *flagValue {
			queryStr += " WHERE tasks." + listCommandConditions[flagKey]
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

	hasRows := false
	for rows.Next() {
		var id int
		var description string
		var completed bool
		var priority string

		err = rows.Scan(&id, &description, &completed, &priority)
		if err != nil {
			return err
		}

		status := "[ ]"
		if completed {
			status = "[o]"
		}
		fmt.Printf("%d.\t %s %s (%s)\n", id, status, description, priority)
		hasRows = true
	}
	if !hasRows {
		fmt.Println("No tasks found")
	}
	return nil
}

func addTask(conn *pgx.Conn, description string, flags map[string]*bool) (int, error) {
	var id, priorityId int = -1, -1
	var err error

	for flagKey, flagValue := range flags {
		if *flagValue {
			priorityId, err = retrievePriorityId(conn, flagKey)
			if err != nil {
				return -1, err
			}
			break
		}
	}

	if priorityId == -1 {
		err = conn.QueryRow(
			context.Background(),
			"INSERT INTO tasks (description) VALUES ($1) RETURNING id",
			description,
		).Scan(&id)
	} else {
		err = conn.QueryRow(
			context.Background(),
			"INSERT INTO tasks (description, priority_id) VALUES ($1, $2) RETURNING id",
			description,
			priorityId,
		).Scan(&id)
	}

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
