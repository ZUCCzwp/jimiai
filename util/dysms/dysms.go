package dysms

import (
	"fmt"
	"jiyu/config"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

func Do(phone string, code int) error {
	conf := sdk.NewConfig()

	credential := credentials.NewAccessKeyCredential(strings.TrimSpace(config.AliyunDypnsConfig.AccessKeyId), strings.TrimSpace(config.AliyunDypnsConfig.AccessKeySecret))
	client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", conf, credential)
	if err != nil {
		return err
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = "基遇"
	request.TemplateCode = "SMS_285780019"
	request.PhoneNumbers = phone
	request.TemplateParam = fmt.Sprintf("{\"code\":\"%d\"}", code)

	result, err := client.SendSms(request)
	fmt.Println(result)

	return err
}
