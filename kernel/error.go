package kernel

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func PanicErrorLang(err error, metadata *fyne.AppMetadata) {
	app.SetMetadata(*metadata)
	a := app.New()
	window := a.NewWindow("GUI for FFmpeg")
	window.SetContent(container.NewVBox(
		widget.NewLabel("Произошла ошибка!"),
		widget.NewLabel("произошла ошибка при получении языковых переводах. \n\r"+err.Error()),
	))
	window.ShowAndRun()
	panic(err.Error())
}
