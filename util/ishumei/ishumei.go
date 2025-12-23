package ishumei

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"jiyu/config"
	"net/http"
)

const (
	EventSign     = "sign"     // 签名
	EventMoment   = "dynamic"  // 动态
	EventNickname = "nickname" // 昵称
	EventComment  = "comment"  // 评论
	EventMessage  = "message"  // 私聊
	EventDefault  = "nickname"
)

// Text 文本检测
func Text(phone, text, eventID string) bool {
	if text == "" {
		return true
	}

	if eventID == "" {
		eventID = EventDefault
	}

	url := "http://api-text-bj.fengkongcloud.com/text/v4"

	h := md5.New()
	h.Write([]byte(phone))
	uid := hex.EncodeToString(h.Sum(nil))

	payload := map[string]interface{}{
		"accessKey": config.ShumeiConfig.AccessKey,
		"appId":     config.ShumeiConfig.AppId,
		"eventId":   eventID,
		"type":      "ALL",
		"data": map[string]string{
			"text":    text,
			"tokenId": uid,
		},
	}
	b, _ := json.Marshal(payload)
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(b))

	var data map[string]interface{}
	if resp != nil {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(respBytes, &data)

		if v, ok := data["riskLevel"]; ok && v == "PASS" {
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}
