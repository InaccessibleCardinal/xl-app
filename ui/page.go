package ui

import (
	"xl-app/xl"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

var fields = []string{"database name", "keywords"}

func run(w fyne.Window, c *fyne.Container) {
	w.SetContent(c)
	w.Resize(fyne.NewSize(800, 700))
	w.ShowAndRun()
}

func makeFileHandler(win fyne.Window) func() {
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

func RenderApp() {
	xlApp := app.New()
	win := xlApp.NewWindow("excel app")

	fileButton := FileUpload(makeFileHandler(win))
	form := CreateForm(fields, func(args ...string) {
		for _, arg := range args {
			println(arg)
		}
	})
	tabs := CreateTabs(TabsConfig{
		"upload file":         fileButton,
		"search the database": form,
	})

	appContainer := container.NewVBox(tabs)
	run(win, appContainer)
}
