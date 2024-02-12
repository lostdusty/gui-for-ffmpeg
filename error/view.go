package error

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/localizer"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ViewContract interface {
	PanicError(err error)
}

type View struct {
	w                fyne.Window
	localizerService localizer.ServiceContract
}

func NewView(w fyne.Window, localizerService localizer.ServiceContract) *View {
	return &View{
		w:                w,
		localizerService: localizerService,
	}
}

func (v View) PanicError(err error) {
	messageHead := v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "error",
	})

	v.w.SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabel(messageHead),
			widget.NewLabel(err.Error()),
		),
		nil,
		nil,
		nil,
		localizer.LanguageSelectionForm(v.localizerService, func(lang localizer.Lang) {
			v.PanicError(err)
		}),
	))
}

func (v View) PanicErrorWriteDirectoryData() {
	message := v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "errorDatabase",
	})
	messageHead := v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "error",
	})

	v.w.SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabel(messageHead),
			widget.NewLabel(message),
		),
		nil,
		nil,
		nil,
		localizer.LanguageSelectionForm(v.localizerService, func(lang localizer.Lang) {
			v.PanicErrorWriteDirectoryData()
		}),
	))
}
