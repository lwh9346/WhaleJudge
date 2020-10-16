package iohelper

import (
	"io"
	"os"
)

//CopyFile 复制文件
func CopyFile(srcFile, destFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(dest, src)
	return err
}
