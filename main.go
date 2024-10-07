package main

import (
	"fmt"
	"image/color"

	"github.com/emarifer/go-fyne-desktop-todoapp/internal/db"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/models"
	"github.com/emarifer/go-fyne-desktop-todoapp/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
	displayText *widget.Label, todos *services.Todos,
) func(di binding.DataItem, co fyne.CanvasObject) {
	return func(di binding.DataItem, co fyne.CanvasObject) {
		t := models.NewTodoFromDataItem(di)
		ctr, _ := co.(*fyne.Container)

		l := ctr.Objects[0].(*widget.Label)
		c := ctr.Objects[1].(*widget.Check)
		ctr.Objects[2].(*widget.Button).OnTapped = func() {
			todos.Remove(t)
			todos.Dbase.DeleteTodo(t)

			fmt.Printf("The ToDo with description %q has been successfully removed!\n", t.Description)
			displayText.SetText(fmt.Sprintf("%q has been successfully removed!", t.Description))
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

func main() {
	// Create and connect to the DB
	db := db.MakeDb()
	defer db.Close()

	// Get data from the DB and bind it to an UntypedList
	todos := services.NewTodosFromDb(&db)
	// defer todos.Persist()

	// Setup App
	a := app.NewWithID("ftodo")
	a.Settings().SetTheme(&forcedVariant{
		Theme:   theme.DefaultTheme(),
		variant: theme.VariantDark,
	})
	w := a.NewWindow("fToDo App")
	w.Resize(fyne.NewSize(480, 600))

	// Keyboard shortcut for closing the application
	ctrlQ := &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}
	w.Canvas().AddShortcut(ctrlQ, func(shortcut fyne.Shortcut) {
		a.Quit()
	})

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

	displayText := &widget.Label{
		Text:       "Display",
		Truncation: fyne.TextTruncateEllipsis,
	}

	deleteBtn := widget.NewButtonWithIcon(
		"Reset", theme.CancelIcon(), func() {
			todos.Drop()

			displayText.SetText("Display")
		},
	)

	list := widget.NewListWithData(
		// the binding.List type
		todos,
		// func that returns the component structure of the List Item
		// exactly the same as the Simple List
		renderListItem,
		// func that is called for each item in the list and allows
		// but this time we get the actual DataItem we need to cast
		bindDataToList(displayText, &todos),
	)
	list.OnSelected = func(id widget.ListItemID) {
		t := todos.All()
		displayText.SetText(t[id].String())
		fmt.Printf("Selected item: %d\n", id)
	}

	w.SetContent(
		container.NewBorder(
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
		),
	)

	w.Canvas().Focus(input)
	w.ShowAndRun()
}

/* REFERENCES:
https://stackoverflow.com/questions/37932551/mkdir-if-not-exists-using-golang

https://stackoverflow.com/questions/71971679/button-action-for-a-specific-list-item-in-fyne

https://stackoverflow.com/questions/66896228/click-event-on-container
https://docs.fyne.io/extend/extending-widgets

Update a collection item given its ID:
https://github.com/ostafen/clover/blob/v2/examples/update/main.go#L32
*/

/* COMMANDS TO BUILD RELEASE:
fyne package --release -exe todoapp
*/
