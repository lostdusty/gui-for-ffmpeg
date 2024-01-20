package helper

import "runtime"

func PathSeparator() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}
