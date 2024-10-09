package context

import (
	"fyne.io/fyne/v2"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/db"
)

type AppContext struct {
	Db db.IDb
	W  fyne.Window
}

func NewAppContext(db db.IDb, w fyne.Window) AppContext {

	return AppContext{
		Db: db,
		W:  w,
	}
}
