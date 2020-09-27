package cloudStorages

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var S3Config map[string]string

func InitAwsS3Config() {
	path_str, _ := filepath.Abs("config/cloudStorages/aws_s3.yml")
	content, err := ioutil.ReadFile(path_str)
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(content, &S3Config)
}
