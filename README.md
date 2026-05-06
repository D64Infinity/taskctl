# taskctl

A simple command-line task manager built with Go and PostgreSQL.

## Features

- Add tasks
- List and count tasks
- Mark tasks as completed
- Delete a task

## Setup

1. Create a PostgreSQL database with the query:
```sql
CREATE DATABASE taskctl;

CREATE TABLE tasks (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    description TEXT NOT NULL,
    completed BOOL NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);
```

2. Create a `.env` file in the project root by copying `.env.example` and update with your credentials (replace `yourpassword` with your actual PostgreSQL password)
3. Build the project:

```bash
go build -o taskctl.exe
```

## Usage

```bash
taskctl list                                # Show all tasks
taskctl list [--pending, -pending]          # Show active tasks
taskctl list [--completed, -completed]      # Show completed tasks 
taskctl count                               # Count all tasks (summary, pending, completed)"
taskctl add "Example task description"      # Add a new task
taskctl done <id>                           # Mark a task as completed
taskctl delete <id>                         # Delete a task
```

## Tech Stack
- Go 1.22+
- PostgreSQL
- pgx driver