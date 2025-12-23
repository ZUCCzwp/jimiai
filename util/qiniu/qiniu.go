package qiniu

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"jiyu/config"
	"log"
	"net/http"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
)

// GenSign 生成签名
//func GenSign(accessKey, secretKey, method, path string, reqParams url.Values) string {
//	reqParams.Set("access_key", accessKey)
//	reqParams.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
//	paramsStr := reqParams.Encode()
//	signStr := fmt.Sprintf("%s\n%s\n%s\n%s", method, path, paramsStr, "")
//	hash := hmac.New(sha1.New, []byte(secretKey))
//	hash.Write([]byte(signStr))
//	signBytes := hash.Sum(nil)
//	sign := base64.URLEncoding.EncodeToString(signBytes)
//	return sign
//}

type TextCensorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Result  struct {
		Scenes struct {
			Antispam struct {
				Details []struct {
					Score float64 `json:"score"`
					Label string  `json:"label"`
				} `json:"details"`
				Suggestion string `json:"suggestion"`
			} `json:"antispam"`
		} `json:"scenes"`
		Suggestion string `json:"suggestion"`
	} `json:"result"`
}

func genSign(accessKey, secretKey string, data string) string {
	mac := qbox.NewMac(accessKey, secretKey)
	return mac.SignWithData([]byte(data))
}

// TextCensor 七牛云文本审核
func TextCensor(text string) (bool, error) {
	accessKey := config.QiniuConfig.AccessKey
	secretKey := config.QiniuConfig.SecretKey

	path := "/v3/text/censor"
	host := "http://ai.qiniuapi.com"
	contentType := "application/json"

	body := make(map[string]any)

	body["data"] = map[string]string{"text": text}
	body["params"] = map[string][]string{"scenes": {"antispam"}}

	jsonStr, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, host+path, bytes.NewBuffer(jsonStr))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", contentType)
	mac := qbox.NewMac(accessKey, secretKey)
	err = mac.AddToken(auth.TokenQiniu, req)

	if err != nil {
		return false, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	if status != 200 {
		return false, errors.New("七牛云文本审核失败:网络错误")
	}

	b, _ := io.ReadAll(resp.Body)
	var result TextCensorResponse
	err = json.Unmarshal(b, &result)
	if err != nil {
		return false, err
	}
	log.Println("qiniu text censor result: ", result)

	return result.Result.Suggestion == "pass", nil
}
