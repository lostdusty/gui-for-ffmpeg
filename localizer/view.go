package localizer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ViewContract interface {
	LanguageSelection(funcSelected func(lang kernel.Lang))
}

type View struct {
	app kernel.AppContract
}

func NewView(app kernel.AppContract) *View {
	return &View{
		app: app,
	}
}

func (v View) LanguageSelection(funcSelected func(lang kernel.Lang)) {
	languages := v.app.GetLocalizerService().GetLanguages()
	listView := widget.NewList(
		func() int {
			return len(languages)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			block := o.(*widget.Label)
			block.SetText(languages[i].Title)
		})
	listView.OnSelected = func(id widget.ListItemID) {
		_ = v.app.GetLocalizerService().SetCurrentLanguage(languages[id])
		funcSelected(languages[id])
	}

	messageHead := v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "languageSelectionHead",
	})
	v.app.GetWindow().SetContent(widget.NewCard(messageHead, "", listView))
}

func LanguageSelectionForm(localizerService kernel.LocalizerContract, funcSelected func(lang kernel.Lang)) fyne.CanvasObject {
	languages := localizerService.GetLanguages()
	currentLanguage := localizerService.GetCurrentLanguage()
	listView := widget.NewList(
		func() int {
			return len(languages)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			block := o.(*widget.Label)
			block.SetText(languages[i].Title)
			if languages[i].Code == currentLanguage.Lang.Code {
				block.TextStyle = fyne.TextStyle{Bold: true}
			}
		})
	listView.OnSelected = func(id widget.ListItemID) {
		_ = localizerService.SetCurrentLanguage(languages[id])
		funcSelected(languages[id])
	}

	messageHead := localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "languageSelectionFormHead",
	})
	return widget.NewCard(messageHead, "", listView)
}
