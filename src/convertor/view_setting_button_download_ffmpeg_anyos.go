//go:build !windows
// +build !windows

package convertor

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (v View) blockDownloadFFmpeg(
	donwloadFFmpeg func(progressBar *widget.ProgressBar, progressMessage *canvas.Text) error,
) *fyne.Container {
	return container.NewVBox()
}
