package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AuthServiceClient struct {
	client  *http.Client
	baseURL string
}

func NewAuthServiceClient(client *http.Client, baseURL string) AuthServiceClient {
	return AuthServiceClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *AuthServiceClient) Register(request RegisterRequest) (*RegisterResponse, int) {
	return makeRequest[RegisterRequest, RegisterResponse](c.client, c.baseURL, "/auth/register", http.MethodPost, &request)
}

func (c *AuthServiceClient) Login(request LoginRequest) (*LoginResponse, int) {
	return makeRequest[LoginRequest, LoginResponse](c.client, c.baseURL, "/auth/login", http.MethodPost, &request)
}

func (c *AuthServiceClient) Refresh(request RefreshRequest) (*RefreshResponse, int) {
	return makeRequest[RefreshRequest, RefreshResponse](c.client, c.baseURL, "/auth/refresh", http.MethodPost, &request)
}

func makeRequest[Req any, Resp any](client *http.Client, baseURL string, endpoint string, method string, request *Req) (*Resp, int) {
	url := fmt.Sprintf("%s%s", baseURL, endpoint)

	var body io.Reader

	if request != nil {
		requestBytes, err := json.Marshal(request)
		if err != nil {
			log.Fatal(err)
		}

		body = bytes.NewBuffer(requestBytes)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode
	}

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response Resp

	if err := json.Unmarshal(responseBytes, &response); err != nil {
		log.Fatal(err)
	}

	return &response, resp.StatusCode
}
