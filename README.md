# taskctl

A simple command-line task manager built with Goland and PostgreSQL.

## Features

- Add tasks
- List all tasks
- Mark tasks as completed

## Setup

1. Create a PostgreSQL database with the query:
CREATE TABLE tasks (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    description TEXT NOT NULL,
    completed BOOL NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

2. Provide a DSN (connection string) in main.go with your PostgreSQL credentials
3. Build the project:

go build -o taskctl.exe

Usage

taskctl list                            # Show all tasks
taskctl add "Example task description"  # Add a new task
taskctl done <id>                       # Mark a task as completed

Tech Stack
- Go 1.22+
- PostgreSQL
- pgx driver