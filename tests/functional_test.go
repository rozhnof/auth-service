package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

type Tokens struct {
	Access  string
	Refresh string
}

func RequireRegister(t *testing.T, email string, password string) Tokens {
	request := RegisterRequest{
		Email:    email,
		Password: password,
	}

	response, statusCode := authServiceClient.Register(request)

	require.NotNil(t, response)
	assert.Equal(t, http.StatusOK, statusCode)

	return AssertLoginSuccess(t, email, password)
}

func AssertLoginSuccess(t *testing.T, email string, password string) Tokens {
	request := LoginRequest{
		Email:    email,
		Password: password,
	}

	response, statusCode := authServiceClient.Login(request)

	require.NotNil(t, response)
	assert.Equal(t, http.StatusOK, statusCode)

	return Tokens{
		Access:  response.AccessToken,
		Refresh: response.RefreshToken,
	}
}

func TestRegister(t *testing.T) {
	if err := SetUp(); err != nil {
		t.Fatal(err)
	}

	request := RegisterRequest{
		Email:    "test.email@gmail.com",
		Password: "test-password",
	}

	response, statusCode := authServiceClient.Register(request)

	require.NotNil(t, response)
	assert.Equal(t, http.StatusOK, statusCode)
	AssertLoginSuccess(t, request.Email, request.Password)
}

func TestSecondRegister(t *testing.T) {
	if err := SetUp(); err != nil {
		t.Fatal(err)
	}

	var (
		email         = "test.email@gmail.com"
		password      = "test-password"
		otherPassword = "other-test-password"
	)

	RequireRegister(t, email, password)

	t.Run("should return status conflict", func(t *testing.T) {
		request := RegisterRequest{
			Email:    email,
			Password: password,
		}

		response, statusCode := authServiceClient.Register(request)

		require.Nil(t, response)
		assert.Equal(t, http.StatusConflict, statusCode)
	})

	t.Run("should return status conflict", func(t *testing.T) {
		request := RegisterRequest{
			Email:    email,
			Password: otherPassword,
		}

		response, statusCode := authServiceClient.Register(request)

		require.Nil(t, response)
		assert.Equal(t, http.StatusConflict, statusCode)
	})
}

func TestLogin(t *testing.T) {
	if err := SetUp(); err != nil {
		t.Fatal(err)
	}

	var (
		email         = "test.email@gmail.com"
		password      = "test-password"
		otherPassword = "other-test-password"
	)

	t.Run("login non existent user should return not found", func(t *testing.T) {
		request := RegisterRequest{
			Email:    email,
			Password: password,
		}

		response, statusCode := authServiceClient.Register(request)

		require.Nil(t, response)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("register and login should return status ok", func(t *testing.T) {
		RequireRegister(t, email, password)
	})

	t.Run("login with other password should return nil response and status ok", func(t *testing.T) {
		request := LoginRequest{
			Email:    email,
			Password: otherPassword,
		}

		response, statusCode := authServiceClient.Login(request)

		require.Nil(t, response)
		assert.Equal(t, http.StatusOK, statusCode)
	})
}

func TestResfresh(t *testing.T) {
	if err := SetUp(); err != nil {
		t.Fatal(err)
	}

	var (
		email    = "test.email@gmail.com"
		password = "test-password"
	)

	tokens := RequireRegister(t, email, password)

	t.Run("refresh with valid token should return new tokens and status ok", func(t *testing.T) {
		request := RefreshRequest{
			RefreshToken: tokens.Refresh,
		}

		response, statusCode := authServiceClient.Refresh(request)

		require.NotNil(t, response)
		assert.Equal(t, http.StatusOK, statusCode)

		tokens.Access = response.AccessToken
		tokens.Refresh = response.RefreshToken
	})

	t.Run("second refresh with returned token should return new tokens and status ok", func(t *testing.T) {
		request := RefreshRequest{
			RefreshToken: tokens.Refresh,
		}

		response, statusCode := authServiceClient.Refresh(request)

		require.NotNil(t, response)
		assert.Equal(t, http.StatusOK, statusCode)
	})
}

func TestResfreshWithInvalidToken(t *testing.T) {
	if err := SetUp(); err != nil {
		t.Fatal(err)
	}

	t.Run("refresh with empty token", func(t *testing.T) {
		request := RefreshRequest{
			RefreshToken: "",
		}

		response, statusCode := authServiceClient.Refresh(request)

		require.Nil(t, response)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("refresh with invalid token", func(t *testing.T) {
		request := RefreshRequest{
			RefreshToken: "invalid-refresh-token",
		}

		response, statusCode := authServiceClient.Refresh(request)

		require.Nil(t, response)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
	})
}
