package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func UiButton(text string, handler func()) *widget.Button {
	button := widget.NewButton(text, handler)
	button.Alignment = widget.ButtonAlign(fyne.TextAlignCenter)
	return button
}
