package tools

import "os"

func Getpwd() string {
	pwd, _ := os.Getwd()
	return pwd
}
