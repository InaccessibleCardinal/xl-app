package ui

import "fyne.io/fyne/v2/widget"

type Fields []string

func CreateForm(fields Fields, onSubmit func(formValues ...string)) *widget.Form {
	var formItems []*widget.FormItem
	var formEntries []*widget.Entry
	for _, field := range fields {
		entry := widget.NewEntry()
		formItems = append(formItems, widget.NewFormItem(field, entry))
		formEntries = append(formEntries, entry)
	}
	return &widget.Form{Items: formItems, OnSubmit: func() {
		var argsToReturn []string
		for _, ent := range formEntries {
			argsToReturn = append(argsToReturn, ent.Text)
			ent.SetText("")
			ent.Refresh()
		}
		onSubmit(argsToReturn...)
	}}
}
