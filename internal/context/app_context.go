package context

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/db"
)

type AppContext struct {
	Db      db.IDb
	Route   binding.String
	Theme   binding.String
	Version string
	w       fyne.Window
}

func NewAppContext(
	db db.IDb, initialRoute AppRoute, initialTheme AppTheme, window fyne.Window,
) AppContext {
	route := initialRoute.String()
	theme := initialTheme.String()

	return AppContext{
		Db:    db,
		Route: binding.BindString(&route),
		Theme: binding.BindString(&theme),
		w:     window,
	}
}

func (ac *AppContext) GetWindow() fyne.Window {

	return ac.w
}

func (ac *AppContext) OnRouteChange(callback func()) {
	ac.Route.AddListener(binding.NewDataListener(callback))
}

func (ac *AppContext) OnThemeChange(callback func()) {
	ac.Theme.AddListener(binding.NewDataListener(callback))
}

func (ac *AppContext) CurrentRoute() AppRoute {
	r, _ := ac.Route.Get()

	return RouteFromString(r)
}

func (ac *AppContext) CurrentTheme() AppTheme {
	t, _ := ac.Theme.Get()

	return ThemeFromString(t)
}

func (ac *AppContext) NavigateTo(route AppRoute) {
	ac.Route.Set(route.String())
}

func (ac *AppContext) ChangeThemeTo(theme AppTheme) {
	ac.Theme.Set(theme.String())
}
