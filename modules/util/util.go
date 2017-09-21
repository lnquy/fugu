package util

import (
	log "github.com/sirupsen/logrus"

	"io"
	"os"
	"path"
	"io/ioutil"
)

func GetWd() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get wd: %v", err)
	}
	return wd
}

func EnsureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0777) // TODO: Permission
	}
}

func EnsureFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fw, err := os.Create(path)
		if err != nil {
			log.Errorf("failed to create file: %v", err)
		}
		fw.Close()
	}
}

func WriteToFile(name, data string) error {
	EnsureDir(path.Dir(name))
	return ioutil.WriteFile(name, []byte(data), 0666)
}

func Copy(src, des string) error {
	srcFile, err := os.OpenFile(src, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	EnsureFile(des)
	destFile, err := os.OpenFile(des, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}
	if err := destFile.Sync(); err != nil {
		return err
	}
	return nil
}

func IsDirEmpty(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}
