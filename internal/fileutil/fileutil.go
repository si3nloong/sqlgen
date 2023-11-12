package fileutil

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/samber/lo"
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

func IsDirEmptyFiles(dir string, excluded ...string) bool {
	var found bool
	// var found = errors.New("file is found")
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() ||
			dir != filepath.Dir(path) ||
			lo.Contains(excluded, filepath.Base(path)) ||
			filepath.Ext(path) != ".go" {
			return nil
		}
		found = true
		return filepath.SkipAll
	}); err != nil {
		return false
	} else {
		return !found
	}
}
