package context

import "strings"

type AppTheme uint8

const (
	light         string = "LIGHT"
	dark          string = "DARK"
	invalid_theme string = "INVALID_THEME"
)

const (
	Light AppTheme = iota
	Dark
)

func getThemesMapping() map[AppTheme]string {

	return map[AppTheme]string{
		Light: light,
		Dark:  dark,
	}
}

func getReverseThemesMapping() map[string]AppTheme {

	return map[string]AppTheme{
		light: Light,
		dark:  Dark,
	}
}

func ThemeFromString(theme string) AppTheme {
	theme = strings.ToUpper(theme)
	mapping := getReverseThemesMapping()
	if value, ok := mapping[theme]; ok {
		return value
	}

	return Dark
}

func (a AppTheme) String() string {
	mapping := getThemesMapping()
	if value, ok := mapping[a]; ok {
		return value
	}

	return invalid_theme
}
