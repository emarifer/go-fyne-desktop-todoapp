package models

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/data/binding"
)

type Todo struct {
	Id          string
	Description string    `clover:"description"`
	Done        bool      `clover:"done"`
	CreatedAt   time.Time `clover:"created_at"`
}

func NewTodo(description string) Todo {
	return Todo{Description: description, Done: false, CreatedAt: time.Now()}
}

func NewTodoFromDataItem(di binding.DataItem) *Todo {
	v, _ := di.(binding.Untyped).Get()
	return v.(*Todo)
}

func (t Todo) String() string {
	return fmt.Sprintf(
		"%s • %s", t.Description, t.CreatedAt.Format(time.RFC822Z),
	)
}

func (t *Todo) MarkAsDone() {
	t.Done = true
}

func (t *Todo) MarkAsToDo() {
	t.Done = false
}

func (t *Todo) ToMap() map[string]interface{} {
	result := map[string]interface{}{}
	result["description"] = t.Description
	result["done"] = t.Done
	result["created_at"] = t.CreatedAt // .Format(time.RFC822Z)

	return result
}
