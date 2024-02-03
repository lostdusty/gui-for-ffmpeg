//go:build !windows
// +build !windows

package handler

import "git.kor-elf.net/kor-elf/gui-for-ffmpeg/convertor"

func getPathsToFF() []convertor.FFPathUtilities {
	return []convertor.FFPathUtilities{{"ffmpeg/bin/ffmpeg", "ffmpeg/bin/ffprobe"}, {"ffmpeg", "ffprobe"}}
}
