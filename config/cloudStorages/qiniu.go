package cloudStorages

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var QiniuConfig map[string]string

type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

func InitQiniuConfig() {
	path_str, _ := filepath.Abs("config/cloudStorages/qiniu.yml")
	content, err := ioutil.ReadFile(path_str)
	if err != nil {
		fmt.Printf("error (%v)", err)
		return
	}
	yaml.Unmarshal(content, &QiniuConfig)
}
