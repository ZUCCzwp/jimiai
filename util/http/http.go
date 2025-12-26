package http

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
)

func Get(url string, token string, result any) (status int, body string, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)

	if err != nil {
		return 500, "", err
	}

	defer resp.Body.Close()
	status = resp.StatusCode
	b, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(b, result)

	return status, string(b), err
}

func Post(url string, data any, token string, result any) (status int, body string, err error) {
	return request("POST", url, data, token, result)
}

func Put(url string, data any, token string, result any) (status int, body string, err error) {
	return request("PUT", url, data, token, result)
}

func Delete(url string, data any, token string, result any) (status int, body string, err error) {
	return request("DELETE", url, data, token, result)
}

func request(method string, url string, data any, token string, result any) (status int, body string, err error) {
	jsonStr, err := json.Marshal(data)

	if err != nil {
		return 500, "", err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500, "", err
	}
	defer resp.Body.Close()

	status = resp.StatusCode
	b, err := io.ReadAll(resp.Body)

	if result != nil {
		err = json.Unmarshal(b, result)
	}

	return status, string(b), err
}

// PostMultipart 发送 multipart/form-data POST 请求
// filePath: 文件路径
// formFields: 表单字段 map[string]string
// token: Bearer token
// result: 响应结果
func PostMultipart(url string, filePath string, formFields map[string]string, token string, result any) (status int, body string, err error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return 500, "", err
	}
	defer file.Close()

	// 创建 multipart writer
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加文件字段
	fileName := filepath.Base(filePath)
	fileField, err := writer.CreateFormFile("input_reference", fileName)
	if err != nil {
		return 500, "", err
	}
	_, err = io.Copy(fileField, file)
	if err != nil {
		return 500, "", err
	}

	// 添加其他表单字段
	for key, value := range formFields {
		err = writer.WriteField(key, value)
		if err != nil {
			return 500, "", err
		}
	}

	// 关闭 writer
	err = writer.Close()
	if err != nil {
		return 500, "", err
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return 500, "", err
	}

	// 设置 Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500, "", err
	}
	defer resp.Body.Close()

	status = resp.StatusCode
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return status, "", err
	}

	if result != nil {
		err = json.Unmarshal(b, result)
	}

	return status, string(b), err
}
