package app

import (
	"fmt"
	"image/color"
	"log"

	"github.com/emarifer/go-fyne-desktop-todoapp/configs"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/db"
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

var themesMap = map[c.AppTheme]fyne.ThemeVariant{
	c.Light: theme.VariantLight,
	c.Dark:  theme.VariantDark,
}

type App struct {
	application     fyne.App
	ctx             *c.AppContext
	isLoggerEnabled bool
	views           map[c.AppRoute]func() *fyne.Container
	window          *fyne.Window
}

func NewApp() App {
	// Setup Application & Window
	a := app.NewWithID(configs.AppId)
	w := a.NewWindow(configs.WindowTitle)
	w.Resize(fyne.NewSize(configs.WindowWidth, configs.WindowHeight))
	w.SetFixedSize(configs.WindowFixed)

	// Keyboard shortcut for closing the application
	ctrlQ := &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}
	w.Canvas().AddShortcut(ctrlQ, func(shortcut fyne.Shortcut) {
		a.Quit()
	})

	// Create and connect to the DB
	db := db.MakeDb(configs.DbName)

	// Setup Context App
	ctx := setupContext(&db, w)
	ctx.Version = configs.Version

	return App{
		application:     a,
		ctx:             &ctx,
		isLoggerEnabled: configs.EnableLogger,
		window:          &w,
		views: map[c.AppRoute]func() *fyne.Container{
			c.List:     func() *fyne.Container { return ui.GetMainView(&ctx) },
			c.Settings: func() *fyne.Container { return ui.GetSettingsView(&ctx) },
		},
	}
}

func (a *App) getView() *fyne.Container {
	key := a.ctx.CurrentRoute()

	if content, ok := a.views[key]; ok {
		return content()
	}

	return a.views[c.List]()
}

func (a *App) setView() {
	(*a.window).SetContent(a.getView())
}

func (a *App) getVariant() fyne.ThemeVariant {
	key := a.ctx.CurrentTheme()

	if variant, ok := themesMap[key]; ok {
		return variant
	}

	return themesMap[c.Dark]
}

func (a *App) setVariant() {
	a.application.Settings().SetTheme(&forcedVariant{
		Theme:   theme.DefaultTheme(),
		variant: a.getVariant(),
	})
}

func (a *App) log(msg string) {
	if a.isLoggerEnabled {
		log.Println(msg)
	}
}

func (a *App) Run() {
	// adding the callback to the listener for view change
	a.ctx.OnRouteChange(func() {
		value := a.ctx.CurrentRoute()
		// log.Printf("route state changed %s", value)
		a.log(fmt.Sprintf("route state changed %s", value))

		a.setView()
	})

	// adding the callback to the listener for changing themes
	a.ctx.OnThemeChange(func() {
		value := a.ctx.CurrentTheme()
		// log.Printf("route state changed %s", value)
		a.log(fmt.Sprintf("theme state changed %s", value))

		a.setVariant()
	})

	a.setVariant()
	a.setView()
	(*a.window).ShowAndRun()

	// log.Println("exiting...")
	a.log("exiting...")
}

func (a *App) Cleanup() {
	// log.Println("Running cleanup")
	a.log("Running cleanup")
	a.ctx.Db.Close()
	// log.Println("Cleanup finished")
	a.log("Cleanup finished")
}

func setupContext(db *db.Db, w fyne.Window) c.AppContext {
	// TODO: in a real application, a condition could be placed here,
	// e.g. the user's login state, to set an initial view in the context.

	return c.NewAppContext(db, configs.InitialRoute, configs.InitialTheme, w)
}
