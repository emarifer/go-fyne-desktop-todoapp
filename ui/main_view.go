package ui

import (
	"fmt"

	"github.com/emarifer/go-fyne-desktop-todoapp/configs"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/models"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	c "github.com/emarifer/go-fyne-desktop-todoapp/internal/context"
)

type tappableEntry struct {
	widget.Entry
}

func newTappableEntry() *tappableEntry {
	e := &tappableEntry{
		widget.Entry{
			PlaceHolder: "Display",
			TextStyle:   fyne.TextStyle{Monospace: true},
		},
	}
	e.ExtendBaseWidget(e)

	return e
}

func (e *tappableEntry) Tapped(_ *fyne.PointEvent) {
	e.Disable()
}

func renderListItem() fyne.CanvasObject {
	return container.NewBorder(
		nil, nil, // Top & bottom
		// ↓ left of the border ↓
		widget.NewCheck("", nil), // func(b bool) {}
		// ↓ right of the border ↓
		widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
		// take the rest of the space ↓
		widget.NewLabel(""),
	)
}

func bindDataToList(
	displayText *tappableEntry, todos *services.Todos, w fyne.Window,
) func(di binding.DataItem, co fyne.CanvasObject) {
	return func(di binding.DataItem, co fyne.CanvasObject) {
		t := models.NewTodoFromDataItem(di)
		ctr, _ := co.(*fyne.Container)

		l := ctr.Objects[0].(*widget.Label)
		c := ctr.Objects[1].(*widget.Check)
		ctr.Objects[2].(*widget.Button).OnTapped = func() {
			msg := fmt.Sprintf("Are you sure you want to delete the task with Description %q?", t.Description)
			dialog.ShowConfirm("Confirmation", msg, func(b bool) {
				if !b {
					return
				}
				todos.Remove(t)
				todos.Dbase.DeleteTodo(t)

				if configs.EnableLogger {
					fmt.Printf("The ToDo with description %q has been successfully removed!\n", t.Description)
				}
				displayText.SetText(fmt.Sprintf("%q has been successfully removed!", t.Description))
			}, w)
		}

		l.Bind(binding.BindString(&t.Description))
		c.Bind(binding.BindBool(&t.Done))

		l.Truncation = fyne.TextTruncateEllipsis
		c.OnChanged = func(b bool) {
			t.Done = b
			todos.Dbase.UpdateTodo(t)
		}
	}
}

func GetMainView(ctx *c.AppContext) *fyne.Container {
	// Get data from the DB and bind it to an UntypedList
	todos := services.NewTodosFromDb(ctx.Db)

	// Setup Widgets
	input := widget.NewEntry()
	input.PlaceHolder = "New TODO description…"
	addBtn := widget.NewButtonWithIcon(
		"Add", theme.DocumentCreateIcon(), func() {
			t := models.NewTodo(input.Text)
			todos.Add(&t)
			input.SetText("")
		},
	)
	addBtn.Disable()
	input.OnChanged = func(s string) {
		// ↓ so that if we delete characters it will be disabled again ↓
		addBtn.Disable()
		if len(s) > 2 {
			addBtn.Enable()
		}
	}

	displayText := newTappableEntry()

	deleteBtn := widget.NewButtonWithIcon(
		"Reset", theme.ViewRefreshIcon(), func() {
			dialog.ShowConfirm(
				"Confirmation",
				"Are you sure you want to delete all the data you have saved? This action is irreversible!!",
				func(b bool) {
					if !b {
						return
					}

					todos.Drop()

					displayText.SetText("Display")
				}, ctx.GetWindow(),
			)
		},
	)

	settingsBtn := navigateBtn(ctx, theme.SettingsIcon(), c.Settings)

	bottomCont := container.NewBorder(nil, nil, nil, settingsBtn, deleteBtn)

	list := widget.NewListWithData(
		// the binding.List type
		todos,
		// func that returns the component structure of the List Item
		// exactly the same as the Simple List
		renderListItem,
		// func that is called for each item in the list and allows
		// but this time we get the actual DataItem we need to cast
		bindDataToList(displayText, &todos, ctx.GetWindow()),
	)
	list.OnSelected = func(id widget.ListItemID) {
		t := todos.All()
		displayText.SetText(t[id].String())
		displayText.Enable()
		if configs.EnableLogger {
			fmt.Printf("Selected item: %d\n", id)
		}
	}

	return container.NewBorder(
		nil, // TOP of the container
		// this will be a the BOTTOM of the container
		container.NewBorder(
			displayText, // TOP
			bottomCont,  // BOTTOM
			nil,         // LEFT
			addBtn,      // RIGHT
			input,       // take the rest of the space ↓
		),
		nil,  // Left
		nil,  // Right
		list, // the rest will take all the rest of the space
	)
}
