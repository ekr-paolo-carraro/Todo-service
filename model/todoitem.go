package model

import (
	"strconv"
)

type TodoItem struct {
	ID       string `json:"sid" binding:"required"`
	UserID   string `json:"userid" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Executed bool   `json:"done"  binding:"required"`
}

func (todo TodoItem) ToString() string {
	return todo.ID + " " + todo.Title + " " + strconv.FormatBool(todo.Executed)
}

type TodoDelegate interface {
	InitData() error
	GetAllItems() ([]TodoItem, error)
	GetTodo(index int) (*TodoItem, error)
	InsertTodo(item TodoItem) (string, error)
	UpdateTodo(item TodoItem) (string, error)
	//DeleteTodo(index int) (string, error)
}
