package handler

import (
	"errors"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/convertor"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/helper"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ConvertorHandlerContract interface {
	MainConvertor()
	FfPathSelection()
	GetFfmpegVersion() (string, error)
	GetFfprobeVersion() (string, error)
}

type ConvertorHandler struct {
	app                 kernel.AppContract
	convertorView       convertor.ViewContract
	convertorRepository convertor.RepositoryContract
}

func NewConvertorHandler(
	app kernel.AppContract,
	convertorView convertor.ViewContract,
	convertorRepository convertor.RepositoryContract,
) *ConvertorHandler {
	return &ConvertorHandler{
		app:                 app,
		convertorView:       convertorView,
		convertorRepository: convertorRepository,
	}
}

func (h ConvertorHandler) MainConvertor() {
	if h.checkingFFPathUtilities() == true {
		h.convertorView.Main(h.runConvert)
		return
	}
	h.convertorView.SelectFFPath("", "", h.saveSettingFFPath, nil, h.downloadFFmpeg)
}

func (h ConvertorHandler) FfPathSelection() {
	ffmpeg, _ := h.convertorRepository.GetPathFfmpeg()
	ffprobe, _ := h.convertorRepository.GetPathFfprobe()
	h.convertorView.SelectFFPath(ffmpeg, ffprobe, h.saveSettingFFPath, h.MainConvertor, h.downloadFFmpeg)
}

func (h ConvertorHandler) GetFfmpegVersion() (string, error) {
	return h.app.GetConvertorService().GetFFmpegVesrion()
}

func (h ConvertorHandler) GetFfprobeVersion() (string, error) {
	return h.app.GetConvertorService().GetFFprobeVersion()
}

func (h ConvertorHandler) runConvert(setting convertor.HandleConvertSetting) {
	h.app.GetQueue().Add(&kernel.ConvertSetting{
		VideoFileInput: setting.VideoFileInput,
		VideoFileOut: kernel.File{
			Path: setting.DirectoryForSave + helper.PathSeparator() + setting.VideoFileInput.Name + ".mp4",
			Name: setting.VideoFileInput.Name,
			Ext:  ".mp4",
		},
		OverwriteOutputFiles: setting.OverwriteOutputFiles,
	})
}

func (h ConvertorHandler) checkingFFPathUtilities() bool {
	if h.checkingFFPath() == true {
		return true
	}

	pathsToFF := getPathsToFF()
	for _, item := range pathsToFF {
		ffmpegChecking, _ := h.app.GetConvertorService().ChangeFFmpegPath(item.FFmpeg)
		if ffmpegChecking == false {
			continue
		}
		ffprobeChecking, _ := h.app.GetConvertorService().ChangeFFprobePath(item.FFprobe)
		if ffprobeChecking == false {
			continue
		}
		_, _ = h.convertorRepository.SavePathFfmpeg(item.FFmpeg)
		_, _ = h.convertorRepository.SavePathFfprobe(item.FFprobe)
		return true
	}

	return false
}

func (h ConvertorHandler) saveSettingFFPath(ffmpegPath string, ffprobePath string) error {
	ffmpegChecking, _ := h.app.GetConvertorService().ChangeFFmpegPath(ffmpegPath)
	if ffmpegChecking == false {
		errorText := h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorFFmpeg",
		})
		return errors.New(errorText)
	}

	ffprobeChecking, _ := h.app.GetConvertorService().ChangeFFprobePath(ffprobePath)
	if ffprobeChecking == false {
		errorText := h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorFFprobe",
		})
		return errors.New(errorText)
	}

	_, _ = h.convertorRepository.SavePathFfmpeg(ffmpegPath)
	_, _ = h.convertorRepository.SavePathFfprobe(ffprobePath)

	h.MainConvertor()

	return nil
}

func (h ConvertorHandler) checkingFFPath() bool {
	_, err := h.app.GetConvertorService().GetFFmpegVesrion()
	if err != nil {
		return false
	}

	_, err = h.app.GetConvertorService().GetFFprobeVersion()
	if err != nil {
		return false
	}

	return true
}
