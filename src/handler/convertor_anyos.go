//go:build !windows
// +build !windows

package handler

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/convertor"
)

func getPathsToFF() []convertor.FFPathUtilities {
	return []convertor.FFPathUtilities{{"ffmpeg/bin/ffmpeg", "ffmpeg/bin/ffprobe"}, {"ffmpeg", "ffprobe"}}
}

func (h ConvertorHandler) downloadFFmpeg(progressBar *widget.ProgressBar, progressMessage *canvas.Text) (err error) {
	return nil
}
