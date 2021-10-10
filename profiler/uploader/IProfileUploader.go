package uploader

import "gitlab.com/dotpe/mindbenders/profiler/entity"

type IProfileUploader interface {
	Upload(bat entity.Batch, serviceName string) error
}
