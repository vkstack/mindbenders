package uploader

import (
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Uploader struct {
	err    error
	s3sess *s3.S3

	prefpath,
	app,
	bucket string
	pathgetter func() string
}

type PathGenerator func() string

func PathGetterWithApp(app string) PathGenerator {
	return func() string {
		return path.Join(os.Getenv("ENV"), app, time.Now().Format("2006-01-02"))
	}
}

func NewUploaderWithBucket(app, bucket string) *S3Uploader {
	uploader := S3Uploader{
		prefpath: path.Join(os.Getenv("ENV"), app),
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
	target = path.Join(sm.prefpath, target)
	return sm.Upload(target, rs)
}

func (sm *S3Uploader) Upload(target string, rs io.ReadSeeker) error {
	_, err := sm.s3sess.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(sm.bucket),
		Key:    aws.String(target),
		Body:   rs,

		ServerSideEncryption: aws.String("AES256"),
		// ContentLength:        aws.Int64(int64(len(fileBlob))),
		// ContentType:          aws.String(contentTypeStr),
	})
	return err
}

// func (sm *S3Uploader) Upload(bat entity.Batch, serviceName string) error {
// 	if sm.s3sess == nil {
// 		return sm.err
// 	}

// 	hostname, err := os.Hostname()
// 	if err != nil {
// 		log.Fatal(err)
// 		hostname = "unknown"
// 	}

// 	// 2006-01-02 15:04:05
// 	target := fmt.Sprintf(
// 		"%s/%s/%s/%s",
// 		os.Getenv("ENV"),
// 		serviceName,
// 		hostname,
// 		bat.Start.Local().Format("2006-01-02/15:04:05"),
// 	)

// 	// for _, prof := range bat.Profiles {
// 	// 	fileDest := fmt.Sprintf("%s.%s", target, prof.Name)
// 	// 	fileBlob := prof.Data
// 	// 	contentTypeStr := http.DetectContentType(fileBlob)
// 	// 	_, err = sm.s3sess.PutObject(&s3.PutObjectInput{
// 	// 		Bucket:               aws.String(S3BUCKET),
// 	// 		Key:                  aws.String(fileDest),
// 	// 		Body:                 bytes.NewReader(fileBlob),
// 	// 		ContentLength:        aws.Int64(int64(len(fileBlob))),
// 	// 		ContentType:          aws.String(contentTypeStr),
// 	// 		ServerSideEncryption: aws.String("AES256"),
// 	// 	})
// 	// 	if err != nil {
// 	// 		log.Println(err)
// 	// 	}
// 	// }
// 	return nil
// }
