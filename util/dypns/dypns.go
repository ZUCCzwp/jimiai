package dypns

import (
	"fmt"
	"jiyu/config"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dypnsapi"
)

// Do 阿里云一键登录
func Do(token string) (string, error) {
	cfg := sdk.NewConfig()

	credential := credentials.NewAccessKeyCredential(strings.TrimSpace(config.AliyunDypnsConfig.AccessKeyId), strings.TrimSpace(config.AliyunDypnsConfig.AccessKeySecret))
	client, err := dypnsapi.NewClientWithOptions("cn-hangzhou", cfg, credential)
	if err != nil {
		log.Println("dypns.Do 阿里云一键登录NewClientWithOptions失败, error:", err)
		return "", err
	}

	request := dypnsapi.CreateGetMobileRequest()
	request.Scheme = "https"
	request.AccessToken = token
	response, err := client.GetMobile(request)
	if err != nil {
		log.Println("dypns.Do 阿里云一键登录GetMobile失败, error:", err)
		return "", err
	}
	fmt.Println(response)
	return response.GetMobileResultDTO.Mobile, nil
}
