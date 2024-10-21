package services

import (
	"time"

	"github.com/emarifer/go-fyne-desktop-todoapp/internal/db"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/models"

	"fyne.io/fyne/v2/data/binding"
)

type Todos struct {
	binding.UntypedList // composition
	Dbase               db.IDb
}

func NewTodosFromDb(db db.IDb) Todos {
	todoList := db.GetAllTodos()

	return newTodos(db, todoList)
}

func newTodos(db db.IDb, todos []models.Todo) Todos {
	t := Todos{
		binding.NewUntypedList(),
		db,
	}

	for _, td := range todos {
		t.Add(&td)
	}

	return t
}

func (t *Todos) Add(todo *models.Todo) {
	// If created_at is the value 'zero' of time.Time,
	// we insert the data into the DB
	var dt *time.Time
	if todo.CreatedAt.String() == "0001-01-01 00:00:00 +0000 UTC" {
		dt, _ = t.Dbase.InsertTodo(todo)
		todo.CreatedAt = *dt
	}

	t.Prepend(todo)
}

func (t *Todos) All() []*models.Todo {
	result := []*models.Todo{}
	for i := 0; i < t.Length(); i++ {
		di, err := t.GetItem(i)
		if err != nil {
			break
		}
		result = append(result, models.NewTodoFromDataItem(di))
	}

	return result
}

func (t *Todos) Drop() {
	t.Dbase.Drop()

	// list, _ := t.Get()
	// list = list[:0]
	t.Set([]any{})
}
