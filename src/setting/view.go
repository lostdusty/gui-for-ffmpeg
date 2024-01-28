package setting

import (
	"ffmpegGui/helper"
	"ffmpegGui/localizer"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"image/color"
	"net/url"
)

type ViewContract interface {
	SelectFFPath(func(ffmpegPath string, ffprobePath string) error)
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

func (v View) SelectFFPath(save func(ffmpegPath string, ffprobePath string) error) {
	errorMessage := canvas.NewText("", color.RGBA{255, 0, 0, 255})
	errorMessage.TextSize = 16
	errorMessage.TextStyle = fyne.TextStyle{Bold: true}

	ffmpegPath, buttonFFmpeg, buttonFFmpegMessage := v.getButtonSelectFile()
	ffprobePath, buttonFFprobe, buttonFFprobeMessage := v.getButtonSelectFile()

	link := widget.NewHyperlink("https://ffmpeg.org/download.html", &url.URL{
		Scheme: "https",
		Host:   "ffmpeg.org",
		Path:   "download.html",
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text: v.localizerService.GetMessage(&i18n.LocalizeConfig{
					MessageID: "titleDownloadLink",
				}),
				Widget: link,
			},
			{
				Text: v.localizerService.GetMessage(&i18n.LocalizeConfig{
					MessageID: "pathToFfmpeg",
				}),
				Widget: buttonFFmpeg,
			},
			{
				Widget: buttonFFmpegMessage,
			},
			{
				Text: v.localizerService.GetMessage(&i18n.LocalizeConfig{
					MessageID: "pathToFfprobe",
				}),
				Widget: buttonFFprobe,
			},
			{
				Widget: buttonFFprobeMessage,
			},
			{
				Widget: errorMessage,
			},
		},
		SubmitText: v.localizerService.GetMessage(&i18n.LocalizeConfig{
			MessageID: "save",
		}),
		OnSubmit: func() {
			err := save(string(*ffmpegPath), string(*ffprobePath))
			if err != nil {
				errorMessage.Text = err.Error()
			}
		},
	}
	selectFFPathTitle := v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "selectFFPathTitle",
	})
	v.w.SetContent(widget.NewCard(selectFFPathTitle, "", container.NewVBox(form)))
}

func (v View) getButtonSelectFile() (filePath *string, button *widget.Button, buttonMessage *canvas.Text) {
	path := ""
	filePath = &path

	buttonMessage = canvas.NewText("", color.RGBA{255, 0, 0, 255})
	buttonMessage.TextSize = 16
	buttonMessage.TextStyle = fyne.TextStyle{Bold: true}

	buttonTitle := v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "choose",
	})

	button = widget.NewButton(buttonTitle, func() {
		fileDialog := dialog.NewFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if err != nil {
					buttonMessage.Text = err.Error()
					setStringErrorStyle(buttonMessage)
					return
				}
				if r == nil {
					return
				}

				path = r.URI().Path()

				buttonMessage.Text = r.URI().Path()
				setStringSuccessStyle(buttonMessage)
			}, v.w)
		helper.FileDialogResize(fileDialog, v.w)
		fileDialog.Show()
	})

	return
}

func setStringErrorStyle(text *canvas.Text) {
	text.Color = color.RGBA{255, 0, 0, 255}
	text.Refresh()
}

func setStringSuccessStyle(text *canvas.Text) {
	text.Color = color.RGBA{49, 127, 114, 255}
	text.Refresh()
}
