package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func GetHash(path string) (hash string) {
	file, err := os.Open(path)
	if err == nil {
		h_ob := sha256.New()
		_, err := io.Copy(h_ob, file)
		if err == nil {
			hash := h_ob.Sum(nil)
			hashvalue := hex.EncodeToString(hash)
			return hashvalue
		} else {
			return "something wrong when use sha256 interface..."
		}
	} else {
		fmt.Printf("failed to open %s\n", path)
	}
	defer file.Close()
	return
}

func GetHash2(path string) (hash string, fileName string) {
	file, err := os.Open(path)
	if err == nil {
		h_ob := sha256.New()
		_, err := io.Copy(h_ob, file)
		if err == nil {
			hash := h_ob.Sum(nil)
			hashvalue := hex.EncodeToString(hash)
			fullName := file.Name()              //路径和文件名包括扩展名
			_, fName := filepath.Split(fullName) //可将路径和文件名分开
			return hashvalue, fName
		} else {
			return "something wrong when use sha256 interface...", "-"
		}
	} else {
		fmt.Printf("failed to open %s\n", path)
	}
	defer file.Close()
	return
}
