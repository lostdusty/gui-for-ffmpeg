package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/convertor"
	error2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/error"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/handler"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
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

var application kernel.AppContract
var ffPathUtilities *kernel.FFPathUtilities

func init() {
	iconResource, _ := fyne.LoadResourceFromPath("icon.png")
	appMetadata := &fyne.AppMetadata{
		ID:      "net.kor-elf.projects.gui-for-ffmpeg",
		Name:    "GUI for FFmpeg",
		Version: "0.4.0",
		Icon:    iconResource,
	}

	localizerService, err := kernel.NewLocalizer("languages", language.Russian)
	if err != nil {
		kernel.PanicErrorLang(err, appMetadata)
	}

	ffPathUtilities = &kernel.FFPathUtilities{FFmpeg: "", FFprobe: ""}
	convertorService := kernel.NewService(ffPathUtilities)
	layoutLocalizerListener := kernel.NewLayoutLocalizerListener()
	localizerService.AddListener(layoutLocalizerListener)

	queue := kernel.NewQueueList()
	application = kernel.NewApp(
		appMetadata,
		localizerService,
		queue,
		kernel.NewQueueLayoutObject(queue, localizerService, layoutLocalizerListener),
		convertorService,
	)
}

func main() {
	errorView := error2.NewView(application)
	if canCreateFile("data/database") != true {
		errorView.PanicErrorWriteDirectoryData()
		application.GetWindow().ShowAndRun()
		return
	}

	db, err := gorm.Open(sqlite.Open("data/database"), &gorm.Config{})
	if err != nil {
		errorView.PanicError(err)
		application.GetWindow().ShowAndRun()
		return
	}

	defer appCloseWithDb(db)

	err = migration.Run(db)
	if err != nil {
		errorView.PanicError(err)
		application.GetWindow().ShowAndRun()
		return
	}

	settingRepository := setting.NewRepository(db)
	convertorRepository := convertor.NewRepository(settingRepository)
	pathFFmpeg, err := convertorRepository.GetPathFfmpeg()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		errorView.PanicError(err)
		application.GetWindow().ShowAndRun()
		return
	}
	ffPathUtilities.FFmpeg = pathFFmpeg

	pathFFprobe, err := convertorRepository.GetPathFfprobe()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		errorView.PanicError(err)
		application.GetWindow().ShowAndRun()
		return
	}
	ffPathUtilities.FFprobe = pathFFprobe

	application.RunConvertor()
	defer application.AfterClosing()

	localizerView := localizer.NewView(application)
	convertorView := convertor.NewView(application)
	convertorHandler := handler.NewConvertorHandler(application, convertorView, convertorRepository)

	localizerRepository := localizer.NewRepository(settingRepository)
	menuView := menu.NewView(application)
	localizerListener := handler.NewLocalizerListener()
	application.GetLocalizerService().AddListener(localizerListener)
	mainMenu := handler.NewMenuHandler(application, convertorHandler, menuView, localizerView, localizerRepository, localizerListener)

	mainHandler := handler.NewMainHandler(application, convertorHandler, mainMenu, localizerRepository)
	mainHandler.Start()

	application.GetWindow().SetMainMenu(mainMenu.GetMainMenu())
	application.GetWindow().ShowAndRun()
}

func appCloseWithDb(db *gorm.DB) {
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
