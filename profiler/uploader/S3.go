package uploader

import (
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Uploader struct {
	err    error
	s3sess *s3.S3

	bucket string
}

func NewUploaderWithBucket(bucket string) *S3Uploader {
	uploader := S3Uploader{
		bucket: bucket,
	}
	s, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	if err != nil {
		uploader.err = err
		log.Fatal(err)
	} else {
		uploader.s3sess = s3.New(s)
	}
	return &uploader
}

func (sm *S3Uploader) UploadProfile(target string, rs io.ReadSeeker) error {
	_, err := sm.s3sess.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(sm.bucket),
		Key:                  aws.String(target),
		Body:                 rs,
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
