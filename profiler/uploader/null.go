package uploader

import "gitlab.com/dotpe/mindbenders/profiler/entity"

type nullUploader struct{}

func NewNullUploader() IProfileUploader {
	return &nullUploader{}
}

func (up *nullUploader) Upload(bat entity.Batch, serviceName string) error {
	return nil
}
