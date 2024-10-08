package ui

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	c "github.com/emarifer/go-fyne-desktop-todoapp/internal/context"
)

func GetSettingsView(ctx *c.AppContext) *fyne.Container {
	url, _ := url.Parse("https://github.com/emarifer/go-fyne-desktop-todoapp")

	// Setup Widgets
	successMsg := newFlashTxtPlaceholder()
	errMsg := newFlashTxtPlaceholder()
	msg := container.NewStack(successMsg, errMsg)

	navigateBackBtn := widget.NewButtonWithIcon(
		"", theme.NavigateBackIcon(), func() {
			ctx.W.SetContent(GetMainView(ctx))
		},
	)

	left := container.NewBorder(nil, navigateBackBtn, nil, nil)

	exportDataBtn := widget.NewButtonWithIcon(
		"Export Data", theme.LogoutIcon(), func() {
			result := ctx.Todos.Dbase.ExportData()
			if result {
				successMessage("Data exported successfully", successMsg)
				return
			}
			errorMessage("The action could not be completed", errMsg)
		},
	)

	importDataBtn := widget.NewButtonWithIcon(
		"Import Data", theme.LoginIcon(), func() {
			// The list needs to be cleared in order to import new data.
			ctx.Todos.Drop()
			result := ctx.Todos.Dbase.ImportData()
			if result {
				successMessage("Data imported successfully", successMsg)
				return
			}
			errorMessage("The action could not be completed!!", errMsg)
		},
	)

	dataManagement := container.NewVBox(importDataBtn, exportDataBtn)

	link := widget.NewHyperlinkWithStyle(
		"https://github.com/emarifer/go-fyne-desktop-todoapp",
		url,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true, Italic: true},
	)

	return container.NewBorder(
		nil, // TOP of the container
		// ↓ this will be a the BOTTOM of the container ↓
		container.NewBorder(nil, nil, left, dataManagement, centered(msg)),
		nil, // Left
		nil, // Right
		container.NewCenter(
			container.NewVBox(
				centered(h1("About")),
				widget.NewLabel(
					"fToDo is a task manager so you don't forget anything 😀",
				),
				small("v1.0.0"),
				&canvas.Text{Text: "", TextSize: 24}, // spacer
				&canvas.Text{
					Text:      "More info:",
					Color:     color.RGBA{207, 130, 37, 255},
					TextSize:  12,
					Alignment: fyne.TextAlignCenter,
					TextStyle: fyne.TextStyle{Bold: true},
				},
				link,
			),
		), // the rest will take all the rest of the space
	)
}
