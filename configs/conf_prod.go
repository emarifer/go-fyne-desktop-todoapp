//go:build prod

package configs

import c "github.com/emarifer/go-fyne-desktop-todoapp/internal/context"

const (
	AppId        = "ftodo_main"
	WindowTitle  = "fToDo App - a mini task manager"
	WindowWidth  = 480
	WindowHeight = 600
	WindowFixed  = true
	InitialRoute = c.List
	InitialTheme = c.Dark
	DbFiles      = "ftodo_db"
	EnableLogger = false
	Version      = "PROD_VERSION"
)
