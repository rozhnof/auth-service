package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const (
	baseURL = "https://www.googleapis.com"
)

type GoogleAuthUserInfo struct {
	Email string `json:"email"`
}

type GoogleAuthClient struct {
	oauthConfig oauth2.Config
}

func NewGoogleAuthClient(config oauth2.Config) *GoogleAuthClient {
	return &GoogleAuthClient{
		oauthConfig: config,
	}
}

func (c *GoogleAuthClient) GetAuthURL(state string) string {
	return c.oauthConfig.AuthCodeURL(state)
}

func (c *GoogleAuthClient) ExchangeOAuthToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token, nil
}

func (c *GoogleAuthClient) GetUserInfo(accessToken string) (*GoogleAuthUserInfo, int, error) {
	const endpoint = "/oauth2/v2/userinfo"

	params := url.Values{}

	params.Add("access_token", accessToken)

	resp, err := http.Get(fmt.Sprintf("%s%s?%s", baseURL, endpoint, params.Encode()))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GoogleAuthUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, resp.StatusCode, nil
}
