package convertor

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/helper"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"image/color"
	"net/url"
)

func (v View) SelectFFPath(
	currentPathFfmpeg string,
	currentPathFfprobe string,
	save func(ffmpegPath string, ffprobePath string) error,
	cancel func(),
) {
	errorMessage := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	errorMessage.TextSize = 16
	errorMessage.TextStyle = fyne.TextStyle{Bold: true}

	ffmpegPath, buttonFFmpeg, buttonFFmpegMessage := v.getButtonSelectFile(currentPathFfmpeg)
	ffprobePath, buttonFFprobe, buttonFFprobeMessage := v.getButtonSelectFile(currentPathFfprobe)

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
			err := save(*ffmpegPath, *ffprobePath)
			if err != nil {
				errorMessage.Text = err.Error()
			}
		},
	}
	if cancel != nil {
		form.OnCancel = cancel
		form.CancelText = v.localizerService.GetMessage(&i18n.LocalizeConfig{
			MessageID: "cancel",
		})
	}
	selectFFPathTitle := v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "selectFFPathTitle",
	})
	v.w.SetContent(widget.NewCard(selectFFPathTitle, "", container.NewVBox(form)))
}

func (v View) getButtonSelectFile(path string) (filePath *string, button *widget.Button, buttonMessage *canvas.Text) {
	filePath = &path

	buttonMessage = canvas.NewText(path, color.RGBA{R: 49, G: 127, B: 114, A: 255})
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
