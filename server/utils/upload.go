package utils

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	"mime/multipart"
	"time"

	"github.com/go-spring/go-spring/spring-boot"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

func init() {
	SpringBoot.RegisterBean(new(UploadService))
}

type OssConfig struct {
	AccessKey string `value:"${oss.access-key}"`
	SecretKey string `value:"${oss.secret-key}"`
	Bucket    string `value:"${oss.bucket}"`
	ImgPath   string `value:"${oss.img-path}"`
}

type UploadService struct {
	OssConfig OssConfig
}

// 接收两个参数 一个文件流 一个 bucket 你的七牛云标准空间的名字
func (service *UploadService) Upload(file *multipart.FileHeader) (err error, path string, key string) {
	putPolicy := storage.PutPolicy{
		Scope: service.OssConfig.Bucket,
	}
	mac := qbox.NewMac(service.OssConfig.AccessKey, service.OssConfig.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	f, e := file.Open()
	if e != nil {
		fmt.Println(e)
		return e, "", ""
	}
	dataLen := file.Size
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	err = formUploader.Put(context.Background(), &ret, upToken, fileKey, f, dataLen, &putExtra)
	if err != nil {
		global.GVA_LOG.Error("upload file fail:", err)
		return err, "", ""
	}
	return err, global.GVA_CONFIG.Qiniu.ImgPath + "/" + ret.Key, ret.Key
}

func (service *UploadService) DeleteFile(key string) error {

	mac := qbox.NewMac(service.OssConfig.AccessKey, service.OssConfig.SecretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	// cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)
	err := bucketManager.Delete(service.OssConfig.Bucket, key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
