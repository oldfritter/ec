package models

import (
	"path"
	"regexp"
	"time"

	"github.com/google/uuid"

	"ec/utils/cloudStorage"
)

type CloudFile struct {
	State       int    `json:"state"`
	StorageName string `json:"storage_name" gorm:"varchar(32)"`
	StorageKey  string `json:"storage_key"`
	OriginName  string `json:"origin_name"`
	FileType    string `json:"file_type"`
}

func (cf *CloudFile) Initialize() {
	cf.OriginName = path.Base(cf.OriginName)
	cf.setStorageName()
	cf.setFileType()
	cf.setStorageKey()
}

func (cf *CloudFile) Attrs() (attrs map[string]string) {
	attrs = make(map[string]string)
	attrs["key"] = cf.StorageKey
	attrs["type"] = cf.FileType
	attrs["name"] = cf.OriginName
	attrs["storage"] = cf.StorageName
	switch cf.StorageName {
	case "s3":
		attrs["url"], _ = cloudStorage.Url4UploadFileToS3(cf.StorageKey)
	case "qiniu":
		attrs["token"] = cloudStorage.Token4UploadFileToQiniu(cf.StorageKey)
	}
	return
}

func (cf *CloudFile) setFileType() {
	if matched, err := regexp.MatchString(".(jpg|jpeg|gif|bmp|png)$", cf.OriginName); matched && err == nil {
		cf.FileType = "image"
	} else if matched, err := regexp.MatchString(".(mp4|avi)$", cf.OriginName); matched && err == nil {
		cf.FileType = "video"
	} else if matched, err := regexp.MatchString(".(mp3)$", cf.OriginName); matched && err == nil {
		cf.FileType = "audio"
	} else {
		cf.FileType = "file"
	}
}

func (cf *CloudFile) setStorageKey() {
	cf.StorageKey = uuid.New().String() + regexp.MustCompile(`\..*`).FindString(cf.OriginName)
}

func (cf *CloudFile) setStorageName() {
	switch time.Now().UnixNano() / 5 {
	case 0:
		cf.StorageName = "s3"
	default:
		cf.StorageName = "qiniu"
	}
}
