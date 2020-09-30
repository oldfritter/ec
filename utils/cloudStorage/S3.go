package cloudStorage

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"ec/config/cloudStorages"
)

func UploadFileToS3(bucket, key, filePath string) error {
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(cloudStorages.S3Config["AWS_REGION"]),
		Credentials: credentials.NewStaticCredentials(cloudStorages.S3Config["AWS_ACCESS_KEY_ID"], cloudStorages.S3Config["AWS_SECRET_ACCESS_KEY"], ""),
	})
	if err != nil {
		return err
	}
	return AddFileToS3(s, bucket, key, filePath)
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func AddFileToS3(s *session.Session, bucket, key, filePath string) error {

	// Open the file for use
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	c, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		ACL:           aws.String("public-read"),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(http.DetectContentType(buffer)),
		// ContentDisposition: aws.String("attachment"),
		// ServerSideEncryption: aws.String("AES256"),
	})
	fmt.Println(c)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func Url4UploadFileToS3(key string) (url string, err error) {
	svc := s3.New(
		session.New(
			&aws.Config{
				Region:      aws.String(cloudStorages.S3Config["AWS_REGION"]),
				Credentials: credentials.NewStaticCredentials(cloudStorages.S3Config["AWS_ACCESS_KEY_ID"], cloudStorages.S3Config["AWS_SECRET_ACCESS_KEY"], ""),
			},
		),
	)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(cloudStorages.S3Config["BUCKET"]),
		Key:    aws.String(key),
	})
	url, err = req.Presign(15 * time.Minute)
	if err != nil {
		log.Println("The err:", err)
	}
	return
}
