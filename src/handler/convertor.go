package handler

import (
	"errors"
	"ffmpegGui/convertor"
	"ffmpegGui/setting"
	"fmt"
	"fyne.io/fyne/v2/widget"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
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
		h.convertorView.Main(h.runConvert, h.getSockPath)
		return
	}
	h.settingView.SelectFFPath(h.saveSettingFFPath)
}

func (h ConvertorHandler) getSockPath(file *convertor.File, progressbar *widget.ProgressBar) (string, error) {
	totalDuration, err := h.getTotalDuration(file)

	if err != nil {
		return "", err
	}
	progressbar.Value = 0
	progressbar.Max = totalDuration
	progressbar.Show()
	progressbar.Refresh()

	rand.Seed(time.Now().Unix())
	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
	l, err := net.Listen("unix", sockFileName)
	if err != nil {
		return "", err
	}

	go func() {
		re := regexp.MustCompile(`frame=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		buf := make([]byte, 16)
		data := ""
		progress := 0.0
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, err := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				if err != nil {
					return
				}
				progress = float64(c)
			}
			if strings.Contains(data, "progress=end") {
				progress = totalDuration
			}
			if progressbar.Value != progress {
				progressbar.Value = progress
				progressbar.Refresh()
			}
		}
	}()

	return sockFileName, nil
}

func (h ConvertorHandler) runConvert(setting convertor.HandleConvertSetting) error {
	return h.convertorService.RunConvert(
		convertor.ConvertSetting{
			VideoFileInput: setting.VideoFileInput,
			SocketPath:     setting.SocketPath,
		},
	)
}

func (h ConvertorHandler) getTotalDuration(file *convertor.File) (float64, error) {
	return h.convertorService.GetTotalDuration(file)
}

func (h ConvertorHandler) checkingFFPathUtilities() bool {
	if h.checkingFFPath() == true {
		return true
	}

	var pathsToFF []convertor.FFPathUtilities
	if runtime.GOOS == "windows" {
		pathsToFF = []convertor.FFPathUtilities{{"ffmpeg/bin/ffmpeg.exe", "ffmpeg/bin/ffprobe.exe"}}
	} else {
		pathsToFF = []convertor.FFPathUtilities{{"ffmpeg/bin/ffmpeg", "ffmpeg/bin/ffprobe"}, {"ffmpeg", "ffprobe"}}
	}
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
		h.settingRepository.Create(ffmpegEntity)
		ffprobeEntity := setting.Setting{Code: "ffprobe", Value: item.FFprobe}
		h.settingRepository.Create(ffprobeEntity)
		return true
	}

	return false
}

func (h ConvertorHandler) saveSettingFFPath(ffmpegPath string, ffprobePath string) error {
	ffmpegChecking, _ := h.convertorService.ChangeFFmpegPath(ffmpegPath)
	if ffmpegChecking == false {
		return errors.New("Это не FFmpeg")
	}

	ffprobeChecking, _ := h.convertorService.ChangeFFprobePath(ffprobePath)
	if ffprobeChecking == false {
		return errors.New("Это не FFprobe")
	}

	ffmpegEntity := setting.Setting{Code: "ffmpeg", Value: ffmpegPath}
	h.settingRepository.Create(ffmpegEntity)
	ffprobeEntity := setting.Setting{Code: "ffprobe", Value: ffprobePath}
	h.settingRepository.Create(ffprobeEntity)

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
