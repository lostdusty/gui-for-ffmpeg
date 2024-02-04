//go:build windows
// +build windows

package handler

import (
	"archive/zip"
	"errors"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/convertor"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func getPathsToFF() []convertor.FFPathUtilities {
	return []convertor.FFPathUtilities{{"ffmpeg\\bin\\ffmpeg.exe", "ffmpeg\\bin\\ffprobe.exe"}}
}

func (h ConvertorHandler) downloadFFmpeg(progressBar *widget.ProgressBar, progressMessage *canvas.Text) (err error) {
	isDirectoryFFmpeg := isDirectory("ffmpeg")
	if isDirectoryFFmpeg == false {
		err = os.Mkdir("ffmpeg", 0777)
		if err != nil {
			return err
		}
	}
	progressMessage.Text = h.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "downloadRun",
	})
	progressMessage.Refresh()
	err = downloadFile("ffmpeg/ffmpeg.zip", "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip", progressBar)
	if err != nil {
		return err
	}

	progressMessage.Text = h.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "unzipRun",
	})
	progressMessage.Refresh()
	err = unZip("ffmpeg/ffmpeg.zip", "ffmpeg")
	if err != nil {
		return err
	}
	_ = os.Remove("ffmpeg/ffmpeg.zip")

	progressMessage.Text = h.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "testFF",
	})
	progressMessage.Refresh()
	err = h.saveSettingFFPath("ffmpeg/ffmpeg-master-latest-win64-gpl/bin/ffmpeg.exe", "ffmpeg/ffmpeg-master-latest-win64-gpl/bin/ffprobe.exe")
	if err != nil {
		return err
	}

	return nil
}

func downloadFile(filepath string, url string, progressBar *widget.ProgressBar) (err error) {
	progressBar.Value = 0
	progressBar.Max = 100

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 32*1024)
	var downloaded int64
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n > 0 {
			f.Write(buf[:n])
			downloaded += int64(n)
			progressBar.Value = float64(downloaded) / float64(resp.ContentLength) * 100
			progressBar.Refresh()
		}
	}
	return nil
}

func unZip(fileZip string, directory string) error {
	archive, err := zip.OpenReader(fileZip)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(directory, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(directory)+string(os.PathSeparator)) {
			return errors.New("invalid file path")
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
