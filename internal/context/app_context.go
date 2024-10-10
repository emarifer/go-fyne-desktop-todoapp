package context

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/db"
)

type AppContext struct {
	Db      db.IDb
	Route   binding.String
	Version string
	w       fyne.Window
}

func NewAppContext(
	db db.IDb, initialRoute AppRoute, window fyne.Window,
) AppContext {
	route := initialRoute.String()

	return AppContext{
		Db:    db,
		Route: binding.BindString(&route),
		w:     window,
	}
}

func (ap *AppContext) GetWindow() fyne.Window {

	return ap.w
}

func (ap *AppContext) OnRouteChange(callback func()) {
	ap.Route.AddListener(binding.NewDataListener(callback))
}

func (ap *AppContext) CurrentRoute() AppRoute {
	r, _ := ap.Route.Get()

	return RouteFromString(r)
}

func (ap *AppContext) NavigateTo(route AppRoute) {
	ap.Route.Set(route.String())
}
