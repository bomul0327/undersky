package s3

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Client struct {
	bucketName string
	sess       *session.Session
}

func NewClient(bucketName string) (*S3Client, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("ap-northeast-2")})
	if err != nil {
		return nil, err
	}

	return &S3Client{
		bucketName: bucketName,
		sess:       sess,
	}, nil
}

func (c *S3Client) Download(key, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(c.sess)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})

	return err
}
