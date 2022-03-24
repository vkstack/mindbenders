package uploader

import "io"

type nullUploader struct{}

func NewNullUploader() IProfileUploader {
	return &nullUploader{}
}

func (up *nullUploader) UploadProfile(target string, rs io.ReadSeeker) error {
	return nil
}
