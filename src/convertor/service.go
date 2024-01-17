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
	GetFFmpegVesrion() (string, error)
	GetFFprobeVersion() (string, error)
	ChangeFFmpegPath(path string) (bool, error)
	ChangeFFprobePath(path string) (bool, error)
}

type FFPathUtilities struct {
	FFmpeg  string
	FFprobe string
}

type Service struct {
	ffPathUtilities *FFPathUtilities
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

func NewService(ffPathUtilities FFPathUtilities) *Service {
	return &Service{
		ffPathUtilities: &ffPathUtilities,
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
	cmd := exec.Command(s.ffPathUtilities.FFmpeg, args...)

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
	cmd := exec.Command(s.ffPathUtilities.FFprobe, args...)
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

func (s Service) GetFFmpegVesrion() (string, error) {
	cmd := exec.Command(s.ffPathUtilities.FFmpeg, "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	text := regexp.MustCompile("\r?\n").Split(strings.TrimSpace(string(out)), -1)
	return text[0], nil
}

func (s Service) GetFFprobeVersion() (string, error) {
	cmd := exec.Command(s.ffPathUtilities.FFprobe, "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	text := regexp.MustCompile("\r?\n").Split(strings.TrimSpace(string(out)), -1)
	return text[0], nil
}

func (s Service) ChangeFFmpegPath(path string) (bool, error) {
	cmd := exec.Command(path, "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	if strings.Contains(strings.TrimSpace(string(out)), "ffmpeg") == false {
		return false, nil
	}
	s.ffPathUtilities.FFmpeg = path
	return true, nil
}

func (s Service) ChangeFFprobePath(path string) (bool, error) {
	cmd := exec.Command(path, "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	if strings.Contains(strings.TrimSpace(string(out)), "ffprobe") == false {
		return false, nil
	}
	s.ffPathUtilities.FFprobe = path
	return true, nil
}
