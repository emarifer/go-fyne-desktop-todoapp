package app

import (
	"image/color"
	"log"

	"github.com/emarifer/go-fyne-desktop-todoapp/internal/db"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/services"
	"github.com/emarifer/go-fyne-desktop-todoapp/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"

	c "github.com/emarifer/go-fyne-desktop-todoapp/internal/context"
)

type forcedVariant struct {
	fyne.Theme

	variant fyne.ThemeVariant
}

func (f *forcedVariant) Color(
	name fyne.ThemeColorName, _ fyne.ThemeVariant,
) color.Color {
	return f.Theme.Color(name, f.variant)
}

type App struct {
	application fyne.App
	ctx         *c.AppContext
	window      *fyne.Window
}

func NewApp() App {
	// Setup Application & Window
	a := app.NewWithID("ftodo")
	a.Settings().SetTheme(&forcedVariant{
		Theme:   theme.DefaultTheme(),
		variant: theme.VariantDark,
	})
	w := a.NewWindow("fToDo App")
	w.Resize(fyne.NewSize(480, 600))
	w.SetFixedSize(true)

	// Keyboard shortcut for closing the application
	ctrlQ := &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}
	w.Canvas().AddShortcut(ctrlQ, func(shortcut fyne.Shortcut) {
		a.Quit()
	})

	// Create and connect to the DB
	db := db.MakeDb()
	// Get data from the DB and bind it to an UntypedList
	todos := services.NewTodosFromDb(&db)

	// Setup Context App
	ctx := c.NewAppContext(todos, w)

	return App{
		application: a,
		ctx:         &ctx,
		window:      &w,
	}
}

func (a *App) setView() {
	(*a.window).SetContent(ui.GetMainView(a.ctx))
}

func (a *App) Run() {
	a.setView()
	(*a.window).ShowAndRun()

	log.Println("exiting...")
}

func (a *App) Cleanup() {
	a.ctx.Todos.Persist()

	log.Println("Running cleanup")
	a.ctx.Todos.Dbase.Close()
	log.Println("Cleanup finished")
}
