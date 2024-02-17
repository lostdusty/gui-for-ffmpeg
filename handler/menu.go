package handler

import (
	"fyne.io/fyne/v2"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/localizer"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/menu"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type MenuHandlerContract interface {
	GetMainMenu() *fyne.MainMenu
	LanguageSelection()
}

type MenuHandler struct {
	app                 kernel.AppContract
	convertorHandler    ConvertorHandlerContract
	menuView            menu.ViewContract
	localizerView       localizer.ViewContract
	localizerRepository localizer.RepositoryContract
	localizerListener   localizerListenerContract
}

func NewMenuHandler(
	app kernel.AppContract,
	convertorHandler ConvertorHandlerContract,
	menuView menu.ViewContract,
	localizerView localizer.ViewContract,
	localizerRepository localizer.RepositoryContract,
	localizerListener localizerListenerContract,
) *MenuHandler {
	return &MenuHandler{
		app:                 app,
		convertorHandler:    convertorHandler,
		menuView:            menuView,
		localizerView:       localizerView,
		localizerRepository: localizerRepository,
		localizerListener:   localizerListener,
	}
}

func (h MenuHandler) GetMainMenu() *fyne.MainMenu {
	settings := h.getMenuSettings()
	help := h.getMenuHelp()

	return fyne.NewMainMenu(settings, help)
}

func (h MenuHandler) getMenuSettings() *fyne.Menu {
	quit := fyne.NewMenuItem(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "exit",
	}), nil)
	quit.IsQuit = true
	h.localizerListener.AddMenuItem("exit", quit)

	languageSelection := fyne.NewMenuItem(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "changeLanguage",
	}), h.LanguageSelection)
	h.localizerListener.AddMenuItem("changeLanguage", languageSelection)

	ffPathSelection := fyne.NewMenuItem(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "changeFFPath",
	}), h.convertorHandler.FfPathSelection)
	h.localizerListener.AddMenuItem("changeFFPath", ffPathSelection)

	settings := fyne.NewMenu(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "settings",
	}), languageSelection, ffPathSelection, quit)
	h.localizerListener.AddMenu("settings", settings)

	return settings
}

func (h MenuHandler) getMenuHelp() *fyne.Menu {
	about := fyne.NewMenuItem(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "about",
	}), h.openAbout)
	h.localizerListener.AddMenuItem("about", about)

	help := fyne.NewMenu(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "help",
	}), about)
	h.localizerListener.AddMenu("help", help)

	return help
}

func (h MenuHandler) openAbout() {
	ffmpeg, err := h.convertorHandler.GetFfmpegVersion()
	if err != nil {
		ffmpeg = h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorFFmpegVersion",
		})
	}
	ffprobe, err := h.convertorHandler.GetFfprobeVersion()
	if err != nil {
		ffprobe = h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorFFprobeVersion",
		})
	}

	h.menuView.About(ffmpeg, ffprobe)
}

func (h MenuHandler) LanguageSelection() {
	h.localizerView.LanguageSelection(func(lang kernel.Lang) {
		_, _ = h.localizerRepository.Save(lang.Code)
		h.convertorHandler.MainConvertor()
	})
}

type menuItems struct {
	menuItem map[string]*fyne.MenuItem
	menu     map[string]*fyne.Menu
}

type LocalizerListener struct {
	menuItems *menuItems
}

type localizerListenerContract interface {
	AddMenu(messageID string, menu *fyne.Menu)
	AddMenuItem(messageID string, menuItem *fyne.MenuItem)
}

func NewLocalizerListener() *LocalizerListener {
	return &LocalizerListener{
		&menuItems{menuItem: map[string]*fyne.MenuItem{}, menu: map[string]*fyne.Menu{}},
	}
}

func (l LocalizerListener) AddMenu(messageID string, menu *fyne.Menu) {
	l.menuItems.menu[messageID] = menu
}

func (l LocalizerListener) AddMenuItem(messageID string, menuItem *fyne.MenuItem) {
	l.menuItems.menuItem[messageID] = menuItem
}

func (l LocalizerListener) Change(localizerService kernel.LocalizerContract) {
	for messageID, menu := range l.menuItems.menuItem {
		menu.Label = localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: messageID})
	}
	for messageID, menu := range l.menuItems.menu {
		menu.Label = localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: messageID})
		menu.Refresh()
	}
}
