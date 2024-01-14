package main

import (
	"ffmpegGui/convertor"
	myError "ffmpegGui/error"
	"ffmpegGui/handler"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const appVersion string = "0.1.0"

func main() {
	a := app.New()
	w := a.NewWindow("GUI FFMpeg!")
	w.Resize(fyne.Size{800, 600})

	errorView := myError.NewView(w)

	pathFFmpeg := "ffmpeg"
	pathFFprobe := "ffprobe"

	convertorView := convertor.NewView(w)
	convertorService := convertor.NewService(pathFFmpeg, pathFFprobe)
	mainHandler := handler.NewConvertorHandler(convertorService, convertorView, errorView)

	mainHandler.GetConvertor()

	w.ShowAndRun()
}
