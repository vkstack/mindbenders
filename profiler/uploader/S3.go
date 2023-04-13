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
	s3sess *s3.S3

	bucket string
}

func NewUploaderWithSession(bucket string, sess *session.Session) *S3Uploader {
	return &S3Uploader{
		bucket: bucket,
		s3sess: s3.New(sess),
	}
}

func NewUploaderWithConfig(bucket string, cfg *aws.Config) *S3Uploader {
	s, err := session.NewSession(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return NewUploaderWithSession(bucket, s)
}

func NewUploaderWithBucket(bucket string) *S3Uploader {
	return NewUploaderWithConfig(bucket, &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
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
