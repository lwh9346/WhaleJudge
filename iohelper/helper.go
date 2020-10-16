package iohelper

import (
	"io"
	"os"
)

//CopyFile 复制文件
func CopyFile(srcFile, destFile string) error {
	src, err := os.Open(srcFile)
	defer src.Close()
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

//WriteStringToFile 将字符串写入文件
func WriteStringToFile(file, s string) error {
	var e error
	_, e = os.Stat(file)
	if e == nil {
		e = os.Remove(file)
		if e != nil {
			return e
		}
	}
	f, e := os.Create(file)
	defer f.Close()
	if e != nil {
		return e
	}
	io.WriteString(f, s)
	return nil
}
