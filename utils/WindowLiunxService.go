package utils

import (
	"runtime"
	"strings"
)

func GetRuntimeDirectory(path string) string {
	systype := runtime.GOOS
	//ForLinux
	if systype == "linux" {
		return GetLinuxDirectory(path)
	}
	//ForWindows
	if systype == "windows" {
		return GetWindowDirectory(path)
	}
	return path
}

func IsWindowRunTime() bool {
	systype := runtime.GOOS
	if systype == "windows" {
		return true
	}
	return false
}

func IsLinuxRunTime() bool {
	systype := runtime.GOOS
	if systype == "linux" {
		return true
	}
	return false
}

func GetLinuxDirectory(path string) string {
	return strings.Replace(path, "\\", "/", -1) //-1表示符合条件的全部替换
}

func GetWindowDirectory(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}

func ReplaceSlashStr(path string) string {
	return strings.Replace(path, "//", "/", -1)
}
