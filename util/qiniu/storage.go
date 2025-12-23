package qiniu

import (
	"context"
	"encoding/json"
	"fmt"
	"jiyu/config"
	"path/filepath"
	"strconv"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

const CdnDomain = "http://cdn1.jiyujiaoyou.cn/"

type ImagePutRet struct {
	Key    string `json:"key"`
	Hash   string `json:"hash"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type ImgPutRet struct {
	Domain string `json:"domain"`
	Key    string `json:"key"`
	Hash   string `json:"hash"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

func PutFile(filePath string) (filename string, width, height int64, err error) {
	accessKey := config.QiniuConfig.AccessKey
	secretKey := config.QiniuConfig.SecretKey

	localFile := filePath
	bucket := config.QiniuConfig.Bucket
	filename = filepath.Base(filePath)
	key := filename

	returnBody := make(map[string]interface{})
	returnBody["key"] = "$(key)"
	returnBody["hash"] = "$(etag)"
	returnBody["width"] = "$(imageInfo.width)"
	returnBody["height"] = "$(imageInfo.height)"
	rb, _ := json.Marshal(returnBody)
	// fmt.Println(string(rb))
	putPolicy := storage.PutPolicy{
		Scope:      bucket,
		ReturnBody: string(rb),
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := ImagePutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			//"x:name": "github logo",
		},
	}
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		// fmt.Println(ret.Key, ret.Hash, ret.Width, ret.Height)
		return "", 0, 0, err
	}

	width, err = strconv.ParseInt(ret.Width, 10, 64)
	if err != nil {
		width = 0
	}
	height, err = strconv.ParseInt(ret.Height, 10, 64)
	if err != nil {
		height = 0
	}

	// fmt.Println(ret.Key, ret.Hash, ret.Width, ret.Height)
	// fmt.Println(ret.Key, ret.Hash)
	return filename, width, height, nil
}

func DeleteFile(filename []string) error {
	return nil
}

func GetImgUploadToken(uid uint, prefix string, fsizeMin, fsizeLimit int64) (upToken string, err error) {
	accessKey := config.QiniuConfig.AccessKey
	secretKey := config.QiniuConfig.SecretKey
	bucket := config.QiniuConfig.Bucket

	flId, err := uuid.NewV4()

	if err != nil {
		return "", err
	}

	flKey := fmt.Sprintf(`%s/%v-%s%s`, prefix, uid, flId.String(), "$(ext)")

	returnBody := ImgPutRet{
		Domain: CdnDomain,
		Key:    "$(key)",
		Hash:   "$(etag)",
		Width:  "$(imageInfo.width)",
		Height: "$(imageInfo.height)",
	}
	rb, _ := json.Marshal(returnBody)

	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
		Scope:        bucket,
		InsertOnly:   1,
		SaveKey:      flKey,
		ForceSaveKey: true,
		ReturnBody:   string(rb),
		FsizeMin:     fsizeMin,
		FsizeLimit:   fsizeLimit,
		MimeLimit:    "image/*",
	}
	return putPolicy.UploadToken(mac), nil
}

func LocalImgUpload(upToken, key, localFile string) (err error) {
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := ImagePutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			//"x:name": "github logo",
		},
	}
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		return err
	}

	return
}
