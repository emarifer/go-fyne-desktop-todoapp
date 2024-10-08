package ui

import (
	"fmt"

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
	displayText *widget.Entry, todos *services.Todos, w fyne.Window,
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

				fmt.Printf("The ToDo with description %q has been successfully removed!\n", t.Description)
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

	// Setup Widgets
	input := widget.NewEntry()
	input.PlaceHolder = "New TODO description…"
	addBtn := widget.NewButtonWithIcon(
		"Add", theme.DocumentCreateIcon(), func() {
			t := models.NewTodo(input.Text)
			ctx.Todos.Add(&t)
			// todos.Add(&t)
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

	displayText := widget.NewEntry()
	displayText.PlaceHolder = "Display"
	displayText.TextStyle = fyne.TextStyle{Monospace: true}
	displayText.Disable()

	deleteBtn := widget.NewButtonWithIcon(
		"Reset", theme.ViewRefreshIcon(), func() {
			dialog.ShowConfirm(
				"Confirmation",
				"Are you sure you want to delete all the data you have saved? This action is irreversible!!",
				func(b bool) {
					if !b {
						return
					}

					ctx.Todos.Drop()

					displayText.SetText("Display")
				}, ctx.W,
			)
		},
	)

	list := widget.NewListWithData(
		// the binding.List type
		ctx.Todos,
		// func that returns the component structure of the List Item
		// exactly the same as the Simple List
		renderListItem,
		// func that is called for each item in the list and allows
		// but this time we get the actual DataItem we need to cast
		bindDataToList(displayText, &ctx.Todos, ctx.W),
	)
	list.OnSelected = func(id widget.ListItemID) {
		t := ctx.Todos.All()
		displayText.SetText(t[id].String())
		fmt.Printf("Selected item: %d\n", id)
	}

	return container.NewBorder(
		nil, // TOP of the container
		// this will be a the BOTTOM of the container
		container.NewBorder(
			displayText, // TOP
			deleteBtn,   // BOTTOM
			nil,         // LEFT
			addBtn,      // RIGHT
			input,       // take the rest of the space ↓
		),
		nil,  // Left
		nil,  // Right
		list, // the rest will take all the rest of the space
	)
}
