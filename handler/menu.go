package handler

import (
	"fyne.io/fyne/v2"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/localizer"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/menu"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type MenuHandlerContract interface {
	GetMainMenu() *fyne.MainMenu
	LanguageSelection()
}

type menuItems struct {
	menuItem map[string]*fyne.MenuItem
	menu     map[string]*fyne.Menu
}

type MenuHandler struct {
	convertorHandler    ConvertorHandlerContract
	menuView            menu.ViewContract
	localizerService    localizer.ServiceContract
	localizerView       localizer.ViewContract
	localizerRepository localizer.RepositoryContract
	menuItems           *menuItems
}

func NewMenuHandler(
	convertorHandler ConvertorHandlerContract,
	menuView menu.ViewContract,
	localizerService localizer.ServiceContract,
	localizerView localizer.ViewContract,
	localizerRepository localizer.RepositoryContract,
) *MenuHandler {
	return &MenuHandler{
		convertorHandler:    convertorHandler,
		menuView:            menuView,
		localizerService:    localizerService,
		localizerView:       localizerView,
		localizerRepository: localizerRepository,
		menuItems:           &menuItems{menuItem: map[string]*fyne.MenuItem{}, menu: map[string]*fyne.Menu{}},
	}
}

func (h MenuHandler) GetMainMenu() *fyne.MainMenu {
	settings := h.getMenuSettings()
	help := h.getMenuHelp()

	return fyne.NewMainMenu(settings, help)
}

func (h MenuHandler) getMenuSettings() *fyne.Menu {
	quit := fyne.NewMenuItem(h.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "exit",
	}), nil)
	quit.IsQuit = true
	h.menuItems.menuItem["exit"] = quit

	languageSelection := fyne.NewMenuItem(h.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "changeLanguage",
	}), h.LanguageSelection)
	h.menuItems.menuItem["changeLanguage"] = languageSelection

	ffPathSelection := fyne.NewMenuItem(h.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "changeFFPath",
	}), h.convertorHandler.FfPathSelection)
	h.menuItems.menuItem["changeFFPath"] = ffPathSelection

	settings := fyne.NewMenu(h.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "settings",
	}), languageSelection, ffPathSelection, quit)
	h.menuItems.menu["settings"] = settings

	return settings
}

func (h MenuHandler) getMenuHelp() *fyne.Menu {
	about := fyne.NewMenuItem(h.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "about",
	}), h.openAbout)
	h.menuItems.menuItem["about"] = about

	help := fyne.NewMenu(h.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "help",
	}), about)
	h.menuItems.menu["help"] = help

	return help
}

func (h MenuHandler) openAbout() {
	ffmpeg, err := h.convertorHandler.GetFfmpegVersion()
	if err != nil {
		ffmpeg = h.localizerService.GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorFFmpegVersion",
		})
	}
	ffprobe, err := h.convertorHandler.GetFfprobeVersion()
	if err != nil {
		ffprobe = h.localizerService.GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorFFprobeVersion",
		})
	}

	h.menuView.About(ffmpeg, ffprobe)
}

func (h MenuHandler) LanguageSelection() {
	h.localizerView.LanguageSelection(func(lang localizer.Lang) {
		_, _ = h.localizerRepository.Save(lang.Code)
		h.menuMessageReload()
		h.convertorHandler.MainConvertor()
	})
}

func (h MenuHandler) menuMessageReload() {
	for messageID, menu := range h.menuItems.menuItem {
		menu.Label = h.localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: messageID})
	}
	for messageID, menu := range h.menuItems.menu {
		menu.Label = h.localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: messageID})
		menu.Refresh()
	}
}
