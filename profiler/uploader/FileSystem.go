package uploader

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

type saveToFileSystem struct{}

func (sm *saveToFileSystem) UploadProfile(target string, rs io.ReadSeeker) error {
	target = filepath.Join("profiles", target)
	if err := os.MkdirAll(path.Dir(target), 0755); err != nil {
		return err
	}
	if f, err := os.Open(target); err != nil {
		return err
	} else {
		_, err = io.Copy(f, rs)
		return err
	}
}

func GetFileSaver() IProfileUploader {
	return &saveToFileSystem{}
}
