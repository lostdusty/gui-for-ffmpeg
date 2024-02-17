//go:build windows
// +build windows

package convertor

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/image/colornames"
	"image/color"
)

func (v View) blockDownloadFFmpeg(
	donwloadFFmpeg func(progressBar *widget.ProgressBar, progressMessage *canvas.Text) error,
) *fyne.Container {

	errorDownloadFFmpegMessage := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	errorDownloadFFmpegMessage.TextSize = 16
	errorDownloadFFmpegMessage.TextStyle = fyne.TextStyle{Bold: true}

	progressDownloadFFmpegMessage := canvas.NewText("", color.RGBA{R: 49, G: 127, B: 114, A: 255})
	progressDownloadFFmpegMessage.TextSize = 16
	progressDownloadFFmpegMessage.TextStyle = fyne.TextStyle{Bold: true}

	progressBar := widget.NewProgressBar()

	var buttonDownloadFFmpeg *widget.Button

	buttonDownloadFFmpeg = widget.NewButton(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "download",
	}), func() {
		buttonDownloadFFmpeg.Disable()

		err := donwloadFFmpeg(progressBar, progressDownloadFFmpegMessage)
		if err != nil {
			errorDownloadFFmpegMessage.Text = err.Error()
		}

		buttonDownloadFFmpeg.Enable()
	})

	downloadFFmpegFromSiteMessage := v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "downloadFFmpegFromSite",
	})

	return container.NewVBox(
		canvas.NewLine(colornames.Darkgreen),
		widget.NewCard(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
			MessageID: "buttonDownloadFFmpeg",
		}), "", container.NewVBox(
			widget.NewRichTextFromMarkdown(
				downloadFFmpegFromSiteMessage+" [https://github.com/BtbN/FFmpeg-Builds/releases](https://github.com/BtbN/FFmpeg-Builds/releases)",
			),
			buttonDownloadFFmpeg,
			errorDownloadFFmpegMessage,
			progressDownloadFFmpegMessage,
			progressBar,
		)),
	)
}
