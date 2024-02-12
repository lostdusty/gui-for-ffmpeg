package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/convertor"
	error2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/error"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/handler"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/localizer"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/menu"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/migration"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/setting"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

const appVersion string = "0.3.1"

func main() {
	a := app.New()
	iconResource, err := fyne.LoadResourceFromPath("icon.png")
	if err == nil {
		a.SetIcon(iconResource)
	}
	w := a.NewWindow("GUI for FFmpeg")
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
	convertorRepository := convertor.NewRepository(settingRepository)
	pathFFmpeg, err := convertorRepository.GetPathFfmpeg()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		errorView.PanicError(err)
		w.ShowAndRun()
		return
	}
	pathFFprobe, err := convertorRepository.GetPathFfprobe()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		errorView.PanicError(err)
		w.ShowAndRun()
		return
	}

	ffPathUtilities := convertor.FFPathUtilities{FFmpeg: pathFFmpeg, FFprobe: pathFFprobe}

	localizerView := localizer.NewView(w, localizerService)
	convertorView := convertor.NewView(w, localizerService)
	convertorService := convertor.NewService(ffPathUtilities)
	defer appCloseWithConvert(convertorService)
	convertorHandler := handler.NewConvertorHandler(convertorService, convertorView, convertorRepository, localizerService)

	localizerRepository := localizer.NewRepository(settingRepository)
	menuView := menu.NewView(w, a, appVersion, localizerService)
	mainMenu := handler.NewMenuHandler(convertorHandler, menuView, localizerService, localizerView, localizerRepository)

	mainHandler := handler.NewMainHandler(convertorHandler, mainMenu, localizerRepository, localizerService)
	mainHandler.Start()

	w.SetMainMenu(mainMenu.GetMainMenu())

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
