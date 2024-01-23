//go:build windows
// +build windows

package handler

func getPathsToFF() []convertor.FFPathUtilities {
	return []convertor.FFPathUtilities{{"ffmpeg\\bin\\ffmpeg.exe", "ffmpeg\\bin\\ffprobe.exe"}}
}
