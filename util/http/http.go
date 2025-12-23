package http

import (
	"bytes"
	"github.com/goccy/go-json"
	"io"
	"net/http"
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
