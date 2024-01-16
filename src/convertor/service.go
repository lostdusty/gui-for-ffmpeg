package convertor

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type ServiceContract interface {
	RunConvert(setting ConvertSetting) error
	GetTotalDuration(file *File) (float64, error)
}

type Service struct {
	pathFFmpeg  string
	pathFFprobe string
}

type File struct {
	Path string
	Name string
	Ext  string
}

type ConvertSetting struct {
	VideoFileInput *File
	SocketPath     string
}

type ConvertData struct {
	totalDuration float64
}

func NewService(pathFFmpeg string, pathFFprobe string) *Service {
	return &Service{
		pathFFmpeg:  pathFFmpeg,
		pathFFprobe: pathFFprobe,
	}
}

func (s Service) RunConvert(setting ConvertSetting) error {
	//args := strings.Split("-report -n -c:v libx264", " ")
	//args := strings.Split("-n -c:v libx264", " ")
	//args = append(args, "-progress", "unix://"+setting.SocketPath, "-i", setting.VideoFileInput.Path, "file-out.mp4")
	//args := "-report -n -i " + setting.VideoFileInput.Path + " -c:v libx264 -progress unix://" + setting.SocketPath + " output-file.mp4"
	//args := "-n -i " + setting.VideoFileInput.Path + " -c:v libx264 -progress unix://" + setting.SocketPath + " output-file.mp4"
	//args := "-y -i " + setting.VideoFileInput.Path + " -c:v libx264 -progress unix://" + setting.SocketPath + " output-file.mp4"
	args := []string{"-y", "-i", setting.VideoFileInput.Path, "-c:v", "libx264", "-progress", "unix://" + setting.SocketPath, "output-file.mp4"}
	cmd := exec.Command(s.pathFFmpeg, args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		errStringArr := regexp.MustCompile("\r?\n").Split(strings.TrimSpace(string(out)), -1)
		if len(errStringArr) > 1 {
			return errors.New(errStringArr[len(errStringArr)-1])
		}
		return err
	}

	return nil
}

func (s Service) GetTotalDuration(file *File) (duration float64, err error) {
	args := []string{"-v", "error", "-select_streams", "v:0", "-count_packets", "-show_entries", "stream=nb_read_packets", "-of", "csv=p=0", file.Path}
	cmd := exec.Command(s.pathFFprobe, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		errString := strings.TrimSpace(string(out))
		if len(errString) > 1 {
			return 0, errors.New(errString)
		}
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
}
