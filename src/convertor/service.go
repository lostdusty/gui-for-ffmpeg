package convertor

import (
	"os/exec"
	"strconv"
	"strings"
)

type ServiceContract interface {
	RunConvert(setting ConvertSetting) error
	GetTotalDuration(file File) (float64, error)
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
	VideoFileInput File
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
	args := "-y -i " + setting.VideoFileInput.Path + " -c:v libx264 -progress unix://" + setting.SocketPath + " output-file.mp4"
	cmd := exec.Command("ffmpeg", strings.Split(args, " ")...)

	//stderr, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		return err
	}

	//scanner := bufio.NewScanner(stderr)
	////scanner.Split(bufio.ScanWords)
	//for scanner.Scan() {
	//	m := scanner.Text()
	//	fmt.Println(m)
	//}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetTotalDuration(file File) (duration float64, err error) {
	args := "-v error -select_streams v:0 -count_packets -show_entries stream=nb_read_packets -of csv=p=0 " + file.Path
	cmd := exec.Command(s.pathFFprobe, strings.Split(args, " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
}
