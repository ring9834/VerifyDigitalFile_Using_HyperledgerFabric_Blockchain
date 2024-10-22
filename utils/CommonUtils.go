package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir //strings.Replace(dir, "\\", "/", -1)
}

var _IconExtension = []string{
	".txt", ".jpg", ".png", ".xls", ".xlsx", ".cer", ".dll", ".rar",
	".aspx", ".cs", ".config", ".xml", ".json", ".cshtml", ".log",
}

func GetExtensionIcon(str string) string {
	if strings.TrimSpace(str) == "" {
		return "Non1.jpg"
	}
	var rlt string = ""
	for _, item := range _IconExtension {
		if item == str {
			rlt = item
			break
		}
	}
	if rlt == "" { //有扩展名但没找到对应的icon
		return "Non1.jpg"
	}
	return strings.Replace(rlt, ".", "", -1) + "1.jpg"
}
