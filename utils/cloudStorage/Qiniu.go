package cloudStorage

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"

	"ec/config/cloudStorages"
)

func UploadFileToQiniu(bucket, key, filePath string) error {
	putPolicy := storage.PutPolicy{
		Scope:      fmt.Sprintf("%s:%s", bucket, key),
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	mac := qbox.NewMac(cloudStorages.QiniuConfig["access_key"], cloudStorages.QiniuConfig["secret_key"])
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := cloudStorages.MyPutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "panama logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, &putExtra)
	if err != nil {
		return err
	}
	exec.Command("sh", "-c", "rm -rf "+filePath).Output()
	return nil
}

func Token4UploadFileToQiniu(key string) string {
	putPolicy := storage.PutPolicy{
		Scope:        cloudStorages.QiniuConfig["bucket"],
		CallbackBody: "key=$(key)&hash=$(etag)&bucket=$(bucket)&fsize=$(fsize)&name=$(x:name)",
	}
	return putPolicy.UploadToken(qbox.NewMac(cloudStorages.QiniuConfig["access_key"], cloudStorages.QiniuConfig["secret_key"]))
}
