package handler

import (
	"bufio"
	"errors"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/convertor"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/helper"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/localizer"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type ConvertorHandlerContract interface {
	MainConvertor()
	FfPathSelection()
	GetFfmpegVersion() (string, error)
	GetFfprobeVersion() (string, error)
}

type ConvertorHandler struct {
	convertorService    convertor.ServiceContract
	convertorView       convertor.ViewContract
	convertorRepository convertor.RepositoryContract
	localizerService    localizer.ServiceContract
}

func NewConvertorHandler(
	convertorService convertor.ServiceContract,
	convertorView convertor.ViewContract,
	convertorRepository convertor.RepositoryContract,
	localizerService localizer.ServiceContract,
) *ConvertorHandler {
	return &ConvertorHandler{
		convertorService:    convertorService,
		convertorView:       convertorView,
		convertorRepository: convertorRepository,
		localizerService:    localizerService,
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
	return h.convertorService.GetFFmpegVesrion()
}

func (h ConvertorHandler) GetFfprobeVersion() (string, error) {
	return h.convertorService.GetFFprobeVersion()
}

func (h ConvertorHandler) runConvert(setting convertor.HandleConvertSetting, progressbar *widget.ProgressBar) error {
	totalDuration, err := h.convertorService.GetTotalDuration(setting.VideoFileInput)
	if err != nil {
		return err
	}
	progress := NewProgress(totalDuration, progressbar, h.localizerService)

	return h.convertorService.RunConvert(
		convertor.ConvertSetting{
			VideoFileInput: setting.VideoFileInput,
			VideoFileOut: &convertor.File{
				Path: setting.DirectoryForSave + helper.PathSeparator() + setting.VideoFileInput.Name + ".mp4",
				Name: setting.VideoFileInput.Name,
				Ext:  ".mp4",
			},
			OverwriteOutputFiles: setting.OverwriteOutputFiles,
		},
		progress,
	)
}

func (h ConvertorHandler) checkingFFPathUtilities() bool {
	if h.checkingFFPath() == true {
		return true
	}

	pathsToFF := getPathsToFF()
	for _, item := range pathsToFF {
		ffmpegChecking, _ := h.convertorService.ChangeFFmpegPath(item.FFmpeg)
		if ffmpegChecking == false {
			continue
		}
		ffprobeChecking, _ := h.convertorService.ChangeFFprobePath(item.FFprobe)
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
	ffmpegChecking, _ := h.convertorService.ChangeFFmpegPath(ffmpegPath)
	if ffmpegChecking == false {
		errorText := h.localizerService.GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorFFmpeg",
		})
		return errors.New(errorText)
	}

	ffprobeChecking, _ := h.convertorService.ChangeFFprobePath(ffprobePath)
	if ffprobeChecking == false {
		errorText := h.localizerService.GetMessage(&i18n.LocalizeConfig{
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
	_, err := h.convertorService.GetFFmpegVesrion()
	if err != nil {
		return false
	}

	_, err = h.convertorService.GetFFprobeVersion()
	if err != nil {
		return false
	}

	return true
}

type Progress struct {
	totalDuration    float64
	progressbar      *widget.ProgressBar
	protocol         string
	localizerService localizer.ServiceContract
}

func NewProgress(totalDuration float64, progressbar *widget.ProgressBar, localizerService localizer.ServiceContract) Progress {
	return Progress{
		totalDuration:    totalDuration,
		progressbar:      progressbar,
		protocol:         "pipe:",
		localizerService: localizerService,
	}
}

func (p Progress) GetProtocole() string {
	return p.protocol
}

func (p Progress) Run(stdOut io.ReadCloser, stdErr io.ReadCloser) error {
	isProcessCompleted := false
	var errorText string

	p.progressbar.Value = 0
	p.progressbar.Max = p.totalDuration
	p.progressbar.Refresh()

	progress := 0.0

	go func() {
		scannerErr := bufio.NewReader(stdErr)
		for {
			line, _, err := scannerErr.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				continue
			}
			data := strings.TrimSpace(string(line))
			errorText = data
		}
	}()

	scannerOut := bufio.NewReader(stdOut)
	for {
		line, _, err := scannerOut.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		data := strings.TrimSpace(string(line))
		if strings.Contains(data, "progress=end") {
			p.progressbar.Value = p.totalDuration
			p.progressbar.Refresh()
			isProcessCompleted = true
			break
		}

		re := regexp.MustCompile(`frame=(\d+)`)
		a := re.FindAllStringSubmatch(data, -1)

		if len(a) > 0 && len(a[len(a)-1]) > 0 {
			c, err := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
			if err != nil {
				continue
			}
			progress = float64(c)
		}
		if p.progressbar.Value != progress {
			p.progressbar.Value = progress
			p.progressbar.Refresh()
		}
	}

	if isProcessCompleted == false {
		if len(errorText) == 0 {
			errorText = p.localizerService.GetMessage(&i18n.LocalizeConfig{
				MessageID: "errorConverter",
			})
		}
		return errors.New(errorText)
	}

	return nil
}
