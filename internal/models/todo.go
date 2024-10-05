package models

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/data/binding"
)

type Todo struct {
	Id          int
	Description string
	Done        bool
	CreatedAt   time.Time
}

func NewTodo(id int, description string) Todo {
	return Todo{id, description, false, time.Now()}
}

func NewTodoFromDataItem(di binding.DataItem) Todo {
	v, _ := di.(binding.Untyped).Get()
	return v.(Todo)
}

func (t Todo) String() string {
	return fmt.Sprintf("%s • %s", t.Description, t.CreatedAt.Format(time.RFC822Z))
}
