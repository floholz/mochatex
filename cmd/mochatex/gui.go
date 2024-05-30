package mochatex

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
)

func Gui(errLog, infoLog *log.Logger) {
	a := app.NewWithID("mochatex")
	window := a.NewWindow("MochaTeX")
	window.Resize(fyne.NewSize(1600, 1000))
	window.SetContent(contentLayout())
	window.ShowAndRun()
}

func contentLayout() *fyne.Container {
	return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("MochaTex"))
}
