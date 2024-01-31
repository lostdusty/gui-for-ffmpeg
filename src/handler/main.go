package handler

import (
	"ffmpegGui/localizer"
)

type MainHandler struct {
	convertorHandler    ConvertorHandlerContract
	menuHandler         MenuHandlerContract
	localizerRepository localizer.RepositoryContract
	localizerService    localizer.ServiceContract
}

func NewMainHandler(
	convertorHandler ConvertorHandlerContract,
	menuHandler MenuHandlerContract,
	localizerRepository localizer.RepositoryContract,
	localizerService localizer.ServiceContract,
) *MainHandler {
	return &MainHandler{
		convertorHandler:    convertorHandler,
		menuHandler:         menuHandler,
		localizerRepository: localizerRepository,
		localizerService:    localizerService,
	}
}

func (h MainHandler) Start() {
	language, err := h.localizerRepository.GetCode()
	if err != nil {
		h.menuHandler.LanguageSelection()
		return
	}
	_ = h.localizerService.SetCurrentLanguageByCode(language)

	h.convertorHandler.MainConvertor()
}
