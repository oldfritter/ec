package cloudStorage

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"

	"ec/config/CloudStorages"
)

func UploadFileToS3(bucket, key, filePath string) error {
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(CloudStorages.S3Config["AWS_REGION"]),
		Credentials: credentials.NewStaticCredentials(CloudStorages.S3Config["AWS_ACCESS_KEY_ID"], CloudStorages.S3Config["AWS_SECRET_ACCESS_KEY"], ""), // token can be left blank for now
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
