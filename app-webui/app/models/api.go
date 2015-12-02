package models

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/pottava/golang-microservices/app-webui/app/logs"
	"github.com/pottava/golang-microservices/app-webui/app/misc"
)

const (
	awsEndpoint = "http://aws"
	dbEndpoint  = "http://dbio"
)

// APIHeader represents API response header
type APIHeader struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// APIResponse represents API response
type APIResponse struct {
	Header APIHeader `json:"header"`
}

func aws(method, target, reqest string, response interface{}) error {
	_, err := request(method, awsEndpoint+target, nil, reqest, response)
	return err
}

func db(method, target, reqest string, response interface{}) error {
	_, err := request(method, dbEndpoint+target, nil, reqest, response)
	return err
}

// HTTP Request
func request(method, endpoint string, headers *map[string]string, reqBody string, resJSON interface{}) (resString string, err error) {
	req, _ := http.NewRequest(method, endpoint, strings.NewReader(reqBody))

	req.Header.Add("Accept-Encoding", "gzip")
	if headers != nil {
		for key, value := range *headers {
			req.Header.Set(key, value)
		}
	}
	// Send HTTP Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.Error.Printf("Could not send a HTTP request. Error: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			logs.Error.Printf("Could not parse gzipped content. Error: %v", err)
			return "", err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	// Parse response Body
	if resJSON != nil {
		err = misc.ReadMBJSON(reader, resJSON, 100) // 100MB
		if err != nil {
			logs.Error.Printf("Could not decode response body as a json. Error: %v", err)
		}
		return "", err
	}
	body, err := misc.ReadMB(reader, 100) // 100MB
	if err != nil {
		logs.Error.Printf("Could not read response body. Error: %v", err)
		return "", err
	}
	return string(body), nil
}
