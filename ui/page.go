package ui

import (
	"xl-app/svc"
	"xl-app/xl"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var fields = []string{"database name", "keywords"}

type App struct {
	service *svc.DynamoService
}

func NewApp(s *svc.DynamoService) *App {
	return &App{
		service: s,
	}
}

func (a *App) run(w fyne.Window, c *fyne.Container) {
	w.SetContent(c)
	w.Resize(fyne.NewSize(800, 700))
	w.ShowAndRun()
}

func (a *App) makeFileHandler(win fyne.Window) func() {
	return func() {
		fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
			fileUri := uc.URI()
			println(fileUri.Path())
			result := xl.ProcessXL(fileUri.Path())
			println("ui receiving results:")
			xl.LogJson(result)

		}, win)
		fd.Show()
	}
}

func SubmitBulkUpload(s string) {
	print(s)
}

func (a *App) RenderApp() {
	xlApp := app.New()
	win := xlApp.NewWindow("excel app")

	fileFormItem := widget.NewFormItem("select a file", FileUpload(a.makeFileHandler(win)))
	fileTextFormItem := widget.NewFormItem("name your database", widget.NewEntry())
	uploadEntries := []*widget.FormItem{fileFormItem, fileTextFormItem}

	form := CreateForm(fields, func(args ...string) {
		values := map[string]string{}
		for i, arg := range args {
			fieldToUpdate := fields[i]
			values[fieldToUpdate] = arg
		}
		a.service.SaveEntity(values)
	})
	tabs := CreateTabs(TabsConfig{
		"upload file with text": &widget.Form{
			Items: uploadEntries,
			OnSubmit: func() {
				SubmitBulkUpload("testing")
			},
		},
		"search the database": form,
	})

	appContainer := container.NewVBox(tabs)
	a.run(win, appContainer)
}
