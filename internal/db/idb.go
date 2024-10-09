package db

import "github.com/emarifer/go-fyne-desktop-todoapp/internal/models"

type IDb interface {
	Close()
	DeleteTodo(todo *models.Todo) bool
	Drop() bool
	ExportData() bool
	GetAllTodos() []models.Todo
	ImportData() bool
	InsertTodo(todo *models.Todo) bool
	UpdateTodo(todo *models.Todo) bool
}
