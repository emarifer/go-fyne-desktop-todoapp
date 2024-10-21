package models

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/google/uuid"
)

type Todo struct {
	Id          string
	Description string
	Done        bool
	CreatedAt   time.Time
}

func NewTodo(description string) Todo {
	return Todo{Id: uuid.NewString(), Description: description}
}

func NewTodoFromDataItem(di binding.DataItem) *Todo {
	v, _ := di.(binding.Untyped).Get()
	return v.(*Todo)
}

func (t Todo) String() string {
	done := "❌"
	if t.Done {
		done = "✅"
	}

	return fmt.Sprintf(
		"%s | %s • %s", done, t.Description, t.CreatedAt.Format(time.DateTime),
	)
}

/* time.DateTime:
https://pkg.go.dev/time#pkg-constants
*/
