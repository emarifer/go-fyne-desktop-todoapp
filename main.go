package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/emarifer/go-fyne-desktop-todoapp/internal/models"
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

/* type tappableLabel struct {
	widget.Label
	extraData interface{}
}

func newTappableLabel() *tappableLabel {
	l := &tappableLabel{}
	l.ExtendBaseWidget(l)

	return l
}

func (t *tappableLabel) Tapped(_ *fyne.PointEvent) {
	// display = fmt.Sprintf("Description: %q • Completed: %t\n", t.Text, t.extraData)
}

func (t *tappableLabel) TappedSecondary(_ *fyne.PointEvent) {}

func (t *tappableLabel) SetExtraData(data any) {
	t.extraData = data
} */

var data = []models.Todo{
	models.NewTodo(0, "Some stuff"),
	models.NewTodo(1, "Some more stuff"),
	models.NewTodo(2, "Some other things"),
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&forcedVariant{
		Theme:   theme.DefaultTheme(),
		variant: theme.VariantDark,
	})
	w := a.NewWindow("fToDo App")

	ctrlQ := &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}
	w.Canvas().AddShortcut(ctrlQ, func(shortcut fyne.Shortcut) {
		a.Quit()
	})

	// t := models.NewTodo("Show this on the window")
	todos := binding.NewUntypedList()
	for _, t := range data {
		todos.Append(t)
	}
	count := 2

	newTodoEntry := widget.NewEntry()
	newTodoEntry.PlaceHolder = "New TODO description…"
	addBtn := widget.NewButtonWithIcon(
		"Add", theme.DocumentCreateIcon(), func() {
			// fmt.Printf("You have typed %q!\n", newTodoEntry.Text)
			count++
			t := models.NewTodo(count, newTodoEntry.Text)
			todos.Append(t)
			data = append(data, t)
			newTodoEntry.SetText("")
		})
	addBtn.Disable()
	newTodoEntry.OnChanged = func(s string) {
		addBtn.Disable()
		if len(s) >= 3 {
			addBtn.Enable()
		}
	}

	displayText := &widget.Label{
		Text:       "Display",
		Truncation: fyne.TextTruncateEllipsis,
	}

	deleteBtn := widget.NewButtonWithIcon(
		"Delete All", theme.CancelIcon(), func() {
			list, _ := todos.Get()
			list = list[:0]
			todos.Set(list)
			data = []models.Todo{}
			count = 0
			displayText.SetText("Display")
		},
	)

	list := widget.NewListWithData(
		// the binding.List type
		todos,
		// func that returns the component structure of the List Item
		// exactly the same as the Simple List
		func() fyne.CanvasObject {
			return container.NewBorder(
				nil, nil,
				// ↓ left of the border ↓
				widget.NewCheck("", func(b bool) {}),
				// ↓ right of the border ↓
				widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
				// take the rest of the space ↓
				widget.NewLabel("template"),
			)
		},
		// func that is called for each item in the list and allows
		// but this time we get the actual DataItem we need to cast
		func(di binding.DataItem, co fyne.CanvasObject) {
			todo := models.NewTodoFromDataItem(di)
			ctr, _ := co.(*fyne.Container)
			// ideally we should check `ok` for each one of those casting
			// but we know that they are those types for sure
			l := ctr.Objects[0].(*widget.Label)
			// l := ctr.Objects[0].(*tappableLabel)
			c := ctr.Objects[1].(*widget.Check)
			ctr.Objects[2].(*widget.Button).OnTapped = func() {
				todos.Remove(todo)
				filtered := []models.Todo{}
				for _, t := range data {
					if t.Id != todo.Id {
						filtered = append(filtered, t)
					}
				}
				data = filtered
				count--
				fmt.Printf("The ToDo with description %q has been successfully removed!\n", todo.Description)
				displayText.SetText(fmt.Sprintf("%q has been successfully removed!", todo.Description))
			}
			/*
				diu, _ := di.(binding.Untyped).Get()
				todo := diu.(models.Todo)
			*/
			l.SetText(todo.Description)
			c.SetChecked(todo.Done)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		displayText.SetText(data[id].String())
		fmt.Printf("Selected item: %d\n", id)
	}

	w.SetContent(
		container.NewBorder(
			nil, // TOP of the container
			// this will be a the BOTTOM of the container
			container.NewBorder(
				displayText, // TOP
				// BOTTOM ↓
				deleteBtn,
				nil, // LEFT
				// RIGHT ↓
				addBtn,
				// take the rest of the space ↓
				newTodoEntry,
			),
			nil, // Left
			nil, // Right
			// ↓ Static list ↓
			/* widget.NewList(
				func() int {
					return len(data)
				},
				func() fyne.CanvasObject {
					// return widget.NewLabel("template")
					return container.NewBorder(
						nil, nil, nil,
						// right of the border
						widget.NewCheck("", func(b bool) {}),
						// take the rest of the space ↓
						widget.NewLabel("template"),
					)
				},
				func(lii widget.ListItemID, co fyne.CanvasObject) {
					// co.(*widget.Label).SetText(data[lii].Description)
					ctr, _ := co.(*fyne.Container)
					// ideally we should check `ok` for each one of those casting
					// but we know that they are those types for sure
					l := ctr.Objects[0].(*widget.Label)
					c := ctr.Objects[1].(*widget.Check)
					l.SetText(data[lii].Description)
					c.SetChecked(data[lii].Done)
				},
			), */
			// the rest will take all the rest of the space
			list,
		),
	)

	w.Resize(fyne.NewSize(300, 400))
	w.ShowAndRun()
}

/* REFERENCES:
https://stackoverflow.com/questions/71971679/button-action-for-a-specific-list-item-in-fyne

https://stackoverflow.com/questions/66896228/click-event-on-container
https://docs.fyne.io/extend/extending-widgets
*/

/* COMMANDS TO BUILD RELEASE:
fyne package --release -exe todoapp
*/
