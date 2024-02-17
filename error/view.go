package error

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/localizer"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ViewContract interface {
	PanicError(err error)
}

type View struct {
	app kernel.AppContract
}

func NewView(app kernel.AppContract) *View {
	return &View{
		app: app,
	}
}

func (v View) PanicError(err error) {
	messageHead := v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "error",
	})

	v.app.GetWindow().SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabel(messageHead),
			widget.NewLabel(err.Error()),
		),
		nil,
		nil,
		nil,
		localizer.LanguageSelectionForm(v.app.GetLocalizerService(), func(lang kernel.Lang) {
			v.PanicError(err)
		}),
	))
}

func (v View) PanicErrorWriteDirectoryData() {
	message := v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "errorDatabase",
	})
	messageHead := v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "error",
	})

	v.app.GetWindow().SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabel(messageHead),
			widget.NewLabel(message),
		),
		nil,
		nil,
		nil,
		localizer.LanguageSelectionForm(v.app.GetLocalizerService(), func(lang kernel.Lang) {
			v.PanicErrorWriteDirectoryData()
		}),
	))
}
