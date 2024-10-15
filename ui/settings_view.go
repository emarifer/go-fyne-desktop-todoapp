package ui

import (
	"fmt"
	"image/color"
	"net/url"

	"github.com/emarifer/go-fyne-desktop-todoapp/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	c "github.com/emarifer/go-fyne-desktop-todoapp/internal/context"
)

func GetSettingsView(ctx *c.AppContext) *fyne.Container {
	// Get data from the DB and bind it to an UntypedList
	todos := services.NewTodosFromDb(ctx.Db)
	url, _ := url.Parse("https://github.com/emarifer/go-fyne-desktop-todoapp")

	// Setup Widgets
	successMsg := newFlashTxtPlaceholder()
	errMsg := newFlashTxtPlaceholder()
	msg := container.NewStack(successMsg, errMsg)

	navigateBackBtn := navigateBtn(ctx, theme.NavigateBackIcon(), c.List)

	left := container.NewBorder(nil, navigateBackBtn, nil, nil)

	themeChangeBtn := themeChangeBtn(ctx, theme.ColorPaletteIcon())

	exportDataBtn := widget.NewButtonWithIcon(
		"Export Data", theme.LogoutIcon(), func() {
			result := todos.Dbase.ExportData()
			if result {
				successMessage("Data exported successfully", successMsg)
				return
			}
			errorMessage("The action could not be completed", errMsg)
		},
	)

	importDataBtn := widget.NewButtonWithIcon(
		"Import Data", theme.LoginIcon(), func() {
			result := todos.Dbase.ImportData()
			if result {
				successMessage("Data imported successfully", successMsg)
				return
			}
			errorMessage("The action could not be completed!!", errMsg)
		},
	)

	settingsManagement := container.NewVBox(
		themeChangeBtn,
		importDataBtn,
		exportDataBtn,
	)

	link := widget.NewHyperlinkWithStyle(
		"https://github.com/emarifer/go-fyne-desktop-todoapp",
		url,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true, Italic: true},
	)

	return container.NewBorder(
		nil, // TOP of the container
		// ↓ this will be a the BOTTOM of the container ↓
		container.NewBorder(nil, nil, left, settingsManagement, centered(msg)),
		nil, // Left
		nil, // Right
		container.NewCenter(
			container.NewVBox(
				centered(h1("About")),
				widget.NewLabel(
					"fToDo is a task manager so you don't forget anything 😀",
				),
				small(fmt.Sprintf("version: %s", ctx.Version)),
				&canvas.Text{Text: "", TextSize: 24}, // spacer
				&canvas.Text{
					Text:      "More info:",
					Color:     color.RGBA{121, 196, 252, 255},
					TextSize:  12,
					Alignment: fyne.TextAlignCenter,
					TextStyle: fyne.TextStyle{Bold: true},
				},
				link,
			),
		), // the rest will take all the rest of the space
	)
}
