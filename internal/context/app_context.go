package context

import (
	"fyne.io/fyne/v2"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/services"
)

type AppContext struct {
	Todos services.Todos
	W     fyne.Window
}

func NewAppContext(t services.Todos, w fyne.Window) AppContext {

	return AppContext{
		Todos: t,
		W:     w,
	}
}
