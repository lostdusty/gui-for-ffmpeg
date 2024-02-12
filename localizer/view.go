package localizer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ViewContract interface {
	LanguageSelection(funcSelected func(lang Lang))
}

type View struct {
	w                fyne.Window
	localizerService ServiceContract
}

func NewView(w fyne.Window, localizerService ServiceContract) *View {
	return &View{
		w:                w,
		localizerService: localizerService,
	}
}

func (v View) LanguageSelection(funcSelected func(lang Lang)) {
	languages := v.localizerService.GetLanguages()
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
		_ = v.localizerService.SetCurrentLanguage(languages[id])
		funcSelected(languages[id])
	}

	messageHead := v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "languageSelectionHead",
	})

	v.w.SetContent(widget.NewCard(messageHead, "", listView))
}

func LanguageSelectionForm(localizerService ServiceContract, funcSelected func(lang Lang)) fyne.CanvasObject {
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
