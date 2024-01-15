package handler

import (
	"ffmpegGui/convertor"
	myError "ffmpegGui/error"
	"fmt"
	"fyne.io/fyne/v2/widget"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ConvertorHandler struct {
	convertorService convertor.ServiceContract
	convertorView    convertor.ViewContract
	errorView        myError.ViewContract
}

func NewConvertorHandler(
	convertorService convertor.ServiceContract,
	convertorView convertor.ViewContract,
	errorView myError.ViewContract,
) *ConvertorHandler {
	return &ConvertorHandler{
		convertorService,
		convertorView,
		errorView,
	}
}

func (h ConvertorHandler) GetConvertor() {
	h.convertorView.Main(h.runConvert, h.getSockPath)
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
