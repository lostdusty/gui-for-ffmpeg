package handler

import (
	"bufio"
	"errors"
	"ffmpegGui/convertor"
	"ffmpegGui/helper"
	"ffmpegGui/setting"
	"fyne.io/fyne/v2/widget"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type ConvertorHandler struct {
	convertorService  convertor.ServiceContract
	convertorView     convertor.ViewContract
	settingView       setting.ViewContract
	settingRepository setting.RepositoryContract
}

func NewConvertorHandler(
	convertorService convertor.ServiceContract,
	convertorView convertor.ViewContract,
	settingView setting.ViewContract,
	settingRepository setting.RepositoryContract,
) *ConvertorHandler {
	return &ConvertorHandler{
		convertorService,
		convertorView,
		settingView,
		settingRepository,
	}
}

func (h ConvertorHandler) GetConvertor() {
	if h.checkingFFPathUtilities() == true {
		h.convertorView.Main(h.runConvert)
		return
	}
	h.settingView.SelectFFPath(h.saveSettingFFPath)
}

func (h ConvertorHandler) runConvert(setting convertor.HandleConvertSetting, progressbar *widget.ProgressBar) error {
	totalDuration, err := h.convertorService.GetTotalDuration(setting.VideoFileInput)
	if err != nil {
		return err
	}
	progress := NewProgress(totalDuration, progressbar)

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
		ffmpegEntity := setting.Setting{Code: "ffmpeg", Value: item.FFmpeg}
		_, _ = h.settingRepository.Create(ffmpegEntity)
		ffprobeEntity := setting.Setting{Code: "ffprobe", Value: item.FFprobe}
		_, _ = h.settingRepository.Create(ffprobeEntity)
		return true
	}

	return false
}

func (h ConvertorHandler) saveSettingFFPath(ffmpegPath string, ffprobePath string) error {
	ffmpegChecking, _ := h.convertorService.ChangeFFmpegPath(ffmpegPath)
	if ffmpegChecking == false {
		return errors.New("это не FFmpeg")
	}

	ffprobeChecking, _ := h.convertorService.ChangeFFprobePath(ffprobePath)
	if ffprobeChecking == false {
		return errors.New("это не FFprobe")
	}

	ffmpegEntity := setting.Setting{Code: "ffmpeg", Value: ffmpegPath}
	_, _ = h.settingRepository.Create(ffmpegEntity)
	ffprobeEntity := setting.Setting{Code: "ffprobe", Value: ffprobePath}
	_, _ = h.settingRepository.Create(ffprobeEntity)

	h.GetConvertor()

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

type progress struct {
	totalDuration float64
	progressbar   *widget.ProgressBar
	protocol      string
}

func NewProgress(totalDuration float64, progressbar *widget.ProgressBar) progress {
	return progress{
		totalDuration: totalDuration,
		progressbar:   progressbar,
		protocol:      "pipe:",
	}
}

func (p progress) GetProtocole() string {
	return p.protocol
}

func (p progress) Run(stdOut io.ReadCloser, stdErr io.ReadCloser) error {
	isProcessCompleted := false
	var errorText string

	p.progressbar.Value = 0
	p.progressbar.Max = p.totalDuration
	p.progressbar.Refresh()

	progress := 0.0

	go func() {
		scannerErr := bufio.NewScanner(stdErr)
		for scannerErr.Scan() {
			errorText = scannerErr.Text()
		}
		if err := scannerErr.Err(); err != nil {
			errorText = err.Error()
		}
	}()

	scannerOut := bufio.NewScanner(stdOut)
	for scannerOut.Scan() {
		data := scannerOut.Text()

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
			errorText = "не смогли отконвертировать видео"
		}
		return errors.New(errorText)
	}

	return nil
}
