package util

import (
	"os"
	"path"
	"strings"
)

// CreateDirIfN 创建文件夹如果文件不存在
func CreateDirIfN(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return false
		} else {
			return true
		}
	}

	return false
}

func CreateFileIfN(filename string) bool {
	CreateDirIfN(path.Dir(filename))

	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(filename)
			if err != nil {
				return false
			}
			defer file.Close()
			return true
		}
	}

	return false
}

func GetFilename(fullFilename string) string {
	filename := path.Base(fullFilename)
	idx := strings.LastIndex(filename, ".")

	return filename[:idx]
}