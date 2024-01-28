package main

import (
	"errors"
	"ffmpegGui/convertor"
	error2 "ffmpegGui/error"
	"ffmpegGui/handler"
	"ffmpegGui/localizer"
	"ffmpegGui/migration"
	"ffmpegGui/setting"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

//const appVersion string = "0.2.0"

func main() {
	a := app.New()
	iconResource, err := fyne.LoadResourceFromPath("icon.png")
	if err == nil {
		a.SetIcon(iconResource)
	}
	w := a.NewWindow("GUI FFMpeg!")
	w.Resize(fyne.Size{Width: 800, Height: 600})
	w.CenterOnScreen()

	localizerService, err := localizer.NewService("languages", language.Russian)
	if err != nil {
		panicErrorLang(w, err)
		w.ShowAndRun()
		return
	}

	errorView := error2.NewView(w, localizerService)
	if canCreateFile("data/database") != true {
		errorView.PanicErrorWriteDirectoryData()
		w.ShowAndRun()
		return
	}

	db, err := gorm.Open(sqlite.Open("data/database"), &gorm.Config{})
	if err != nil {
		errorView.PanicError(err)
		w.ShowAndRun()
		return
	}

	defer appCloseWithDb(db)

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

	localizerView := localizer.NewView(w, localizerService)
	convertorView := convertor.NewView(w, localizerService)
	settingView := setting.NewView(w, localizerService)
	convertorService := convertor.NewService(ffPathUtilities)
	defer appCloseWithConvert(convertorService)
	mainHandler := handler.NewConvertorHandler(convertorService, convertorView, settingView, localizerView, settingRepository, localizerService)

	mainHandler.LanguageSelection()

	w.ShowAndRun()
}

func appCloseWithDb(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
}

func appCloseWithConvert(convertorService convertor.ServiceContract) {
	for _, cmd := range convertorService.GetRunningProcesses() {
		_ = cmd.Process.Kill()
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

func panicErrorLang(w fyne.Window, err error) {
	w.SetContent(container.NewVBox(
		widget.NewLabel("Произошла ошибка!"),
		widget.NewLabel("произошла ошибка при получении языковых переводах. \n\r"+err.Error()),
	))
}
