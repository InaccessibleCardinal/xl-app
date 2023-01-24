package ui

import (
	"fyne.io/fyne/v2/widget"
)

func FileUpload(callback func()) *widget.Button {
	fileButton := widget.NewButton("Open file", callback)

	return fileButton
}
