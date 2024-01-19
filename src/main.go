package main

import (
	"errors"
	"ffmpegGui/convertor"
	myError "ffmpegGui/error"
	"ffmpegGui/handler"
	"ffmpegGui/migration"
	"ffmpegGui/setting"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

//const appVersion string = "0.1.1"

func main() {
	a := app.New()
	w := a.NewWindow("GUI FFMpeg!")
	w.Resize(fyne.Size{Width: 800, Height: 600})
	w.CenterOnScreen()

	errorView := myError.NewView(w)

	if canCreateFile("data/database") != true {
		errorView.PanicError(errors.New("не смогли создать файл 'database' в папке 'data'"))
		w.ShowAndRun()
		return
	}

	db, err := gorm.Open(sqlite.Open("data/database"), &gorm.Config{})
	if err != nil {
		errorView.PanicError(err)
		w.ShowAndRun()
		return
	}

	defer appClose(db)

	err = migration.Run(db)
	if err != nil {
		errorView.PanicError(err)
		w.ShowAndRun()
		return
	}

	settingRepository := setting.NewRepository(db)
	pathFFmpeg, err := settingRepository.GetValue("ffmpeg")
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		errorView.PanicError(err)
		w.ShowAndRun()
		return
	}
	pathFFprobe, err := settingRepository.GetValue("ffprobe")
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		errorView.PanicError(err)
		w.ShowAndRun()
		return
	}

	ffPathUtilities := convertor.FFPathUtilities{FFmpeg: pathFFmpeg, FFprobe: pathFFprobe}

	convertorView := convertor.NewView(w)
	settingView := setting.NewView(w)
	convertorService := convertor.NewService(ffPathUtilities)
	mainHandler := handler.NewConvertorHandler(convertorService, convertorView, settingView, settingRepository)

	mainHandler.GetConvertor()

	w.ShowAndRun()
}

func appClose(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
}

func canCreateFile(path string) bool {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return false
	}
	_ = file.Close()
	return true
}
