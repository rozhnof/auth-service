package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

var client = &http.Client{}

func Register(request RegisterRequest, queryParams ...QueryParams) (*RegisterResponse, int, error) {
	return doPost[RegisterRequest, RegisterResponse]("/auth/register", request, "", queryParams)
}

func Login(request LoginRequest, queryParams ...QueryParams) (*LoginResponse, int, error) {
	return doPost[LoginRequest, LoginResponse]("/auth/login", request, "", queryParams)
}

func Refresh(request RefreshRequest, queryParams ...QueryParams) (*RefreshResponse, int, error) {
	return doPost[RefreshRequest, RefreshResponse]("/auth/refresh", request, "", queryParams)
}

func doPost[Req any, Resp any](endpoint string, request Req, accessToken string, queryParams []QueryParams) (*Resp, int, error) {
	const method = "POST"

	url := buildURL(endpoint, queryParams...)

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, nil
	}

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var response Resp

	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, 0, err
	}

	return &response, resp.StatusCode, nil
}

func buildURL(endpoint string, queryParams ...QueryParams) string {
	urlBuilder := strings.Builder{}
	urlBuilder.WriteString(baseURL)
	urlBuilder.WriteString(endpoint)

	for i, queryParam := range queryParams {
		if i == 0 {
			urlBuilder.WriteString("?")
		} else {
			urlBuilder.WriteString("&")
		}

		urlBuilder.WriteString(queryParam.Key)
		urlBuilder.WriteString("=")
		urlBuilder.WriteString(queryParam.Value)
	}

	return urlBuilder.String()
}
