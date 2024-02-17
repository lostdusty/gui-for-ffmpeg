package handler

import (
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/localizer"
)

type MainHandler struct {
	app                 kernel.AppContract
	convertorHandler    ConvertorHandlerContract
	menuHandler         MenuHandlerContract
	localizerRepository localizer.RepositoryContract
}

func NewMainHandler(
	app kernel.AppContract,
	convertorHandler ConvertorHandlerContract,
	menuHandler MenuHandlerContract,
	localizerRepository localizer.RepositoryContract,
) *MainHandler {
	return &MainHandler{
		app:                 app,
		convertorHandler:    convertorHandler,
		menuHandler:         menuHandler,
		localizerRepository: localizerRepository,
	}
}

func (h MainHandler) Start() {
	language, err := h.localizerRepository.GetCode()
	if err != nil {
		h.menuHandler.LanguageSelection()
		return
	}
	_ = h.app.GetLocalizerService().SetCurrentLanguageByCode(language)

	h.convertorHandler.MainConvertor()
}
