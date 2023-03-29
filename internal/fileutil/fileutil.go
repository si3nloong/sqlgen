package fileutil

import (
	"os"
	"path/filepath"
	"runtime"
)

func Getpwd() string {
	pwd, _ := os.Getwd()
	return pwd
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CurDir() string {
	_, path, _, _ := runtime.Caller(1)
	return filepath.Dir(path)
}
