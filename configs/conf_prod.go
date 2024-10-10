//go:build prod

package configs

const (
	AppId        = "ftodo_main"
	WindowTitle  = "fToDo App - a mini task manager"
	WindowWidth  = 480
	WindowHeight = 600
	WindowFixed  = true
	DbFiles      = "ftodo_db"
	EnableLogger = false
	Version      = "PROD_VERSION"
)
