package error

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ViewContract interface {
	PanicError(err error)
}

type View struct {
	w fyne.Window
}

func NewView(w fyne.Window) *View {
	return &View{w}
}

func (v View) PanicError(err error) {
	v.w.SetContent(container.NewVBox(
		widget.NewLabel("Произошла ошибка!"),
		widget.NewLabel("Ошибка: "+err.Error()),
	))
}
