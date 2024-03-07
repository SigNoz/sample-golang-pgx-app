package main

import (
	"context"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func getAllTodos(ctx context.Context) ([]Todo, error) {
	var todos []Todo
	rows, err := conn.Query(ctx, "SELECT id, title, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

// CreateTodo inserts a new todo into the database
func createTodo(ctx context.Context, title string) (Todo, error) {
	var todo Todo
	err := conn.QueryRow(ctx, "INSERT INTO todos (title) VALUES ($1) RETURNING id, title, completed", title).Scan(&todo.ID, &todo.Title, &todo.Completed)
	if err != nil {
		return todo, err
	}
	return todo, nil
}

// UpdateTodo updates a todo's title and completed status
func updateTodo(ctx context.Context, id int, title string, completed bool) (Todo, error) {
	var todo Todo
	err := conn.QueryRow(ctx, "UPDATE todos SET title = $2, completed = $3 WHERE id = $1 RETURNING id, title, completed", id, title, completed).Scan(&todo.ID, &todo.Title, &todo.Completed)
	if err != nil {
		return todo, err
	}
	return todo, nil
}

// DeleteTodo removes a todo from the database
func deleteTodo(ctx context.Context, id int) error {
	_, err := conn.Exec(ctx, "DELETE FROM todos WHERE id = $1", id)
	return err
}
