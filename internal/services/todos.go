package services

import (
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/db"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/models"

	"fyne.io/fyne/v2/data/binding"
)

type Todos struct {
	binding.UntypedList
	Dbase *db.Db
}

func NewTodosFromDb(db *db.Db) Todos {
	todoList := db.GetAllTodos()

	return newTodos(db, todoList)
}

func newTodos(db *db.Db, todos []models.Todo) Todos {
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
	if todo.Id == "" {
		t.Dbase.InsertTodo(todo)
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

	list, _ := t.Get()
	list = list[:0]
	t.Set(list)
}

/* func (t *Todos) Persist() {
	t.Dbase.Save(t.All())
} */
