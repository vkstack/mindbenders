package uploader

import (
	"io"
)

type IProfileUploader interface {
	UploadProfile(target string, rs io.ReadSeeker) error
}
