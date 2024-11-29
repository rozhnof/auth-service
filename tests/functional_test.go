package tests

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/rozhnof/auth-service/internal/infrastructure/database/postgres"
	"github.com/rozhnof/auth-service/internal/pkg/config"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

const (
	testConfigPath = "../config/test-config.yaml"
	baseURL        = "http://localhost:9090"
)

var db postgres.Database

func init() {
	cfg, err := config.NewConfig[config.Postgres](testConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	postgresCfg := postgres.DatabaseConfig{
		Address:  cfg.Address,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
		SSL:      cfg.SSL,
	}

	database, err := postgres.NewDatabase(context.Background(), postgresCfg)
	if err != nil {
		log.Fatal(err)
	}

	db = database
}

func Test1(t *testing.T) {
	if err := truncateAllTables(); err != nil {
		t.Fatal(err)
	}

	t.Run("first register", func(t *testing.T) {
		request := RegisterRequest{
			Email:    "test.email@gmail.com",
			Password: "test-password",
		}

		response, statusCode, err := Register(request)
		if err != nil {
			t.Fatal(err)
		}

		require.NotNil(t, response)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("second register", func(t *testing.T) {
		request := RegisterRequest{
			Email:    "test.email@gmail.com",
			Password: "test-password",
		}

		response, statusCode, err := Register(request)
		if err != nil {
			t.Fatal(err)
		}

		require.Nil(t, response)
		assert.Equal(t, http.StatusConflict, statusCode)
	})

	t.Run("third register with other password", func(t *testing.T) {
		request := RegisterRequest{
			Email:    "test.email@gmail.com",
			Password: "other-test-password",
		}

		response, statusCode, err := Register(request)
		if err != nil {
			t.Fatal(err)
		}

		require.Nil(t, response)
		assert.Equal(t, http.StatusConflict, statusCode)
	})
}

func Test2(t *testing.T) {
	if err := truncateAllTables(); err != nil {
		t.Fatal(err)
	}

	t.Run("register", func(t *testing.T) {
		request := RegisterRequest{
			Email:    "test.email@gmail.com",
			Password: "test-password",
		}

		response, statusCode, err := Register(request)
		if err != nil {
			t.Fatal(err)
		}

		require.NotNil(t, response)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("login", func(t *testing.T) {
		request := LoginRequest{
			Email:    "test.email@gmail.com",
			Password: "test-password",
		}

		response, statusCode, err := Login(request)
		if err != nil {
			t.Fatal(err)
		}

		require.NotNil(t, response)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("update tokens with empty token", func(t *testing.T) {
		request := RefreshRequest{
			RefreshToken: "",
		}

		response, statusCode, err := Refresh(request)
		if err != nil {
			t.Fatal(err)
		}

		require.Nil(t, response)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("update tokens with empty token", func(t *testing.T) {
		request := RefreshRequest{
			RefreshToken: "invalid-refresh-token",
		}

		response, statusCode, err := Refresh(request)
		if err != nil {
			t.Fatal(err)
		}

		require.Nil(t, response)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
	})
}

func Test3(t *testing.T) {
	if err := truncateAllTables(); err != nil {
		t.Fatal(err)
	}

	t.Run("register", func(t *testing.T) {
		request := RegisterRequest{
			Email:    "test.email@gmail.com",
			Password: "test-password",
		}

		response, statusCode, err := Register(request)
		if err != nil {
			t.Fatal(err)
		}

		require.NotNil(t, response)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	var refreshToken1 string

	t.Run("login", func(t *testing.T) {
		request := LoginRequest{
			Email:    "test.email@gmail.com",
			Password: "test-password",
		}

		response, statusCode, err := Login(request)
		if err != nil {
			t.Fatal(err)
		}

		require.NotNil(t, response)
		assert.Equal(t, http.StatusOK, statusCode)

		refreshToken1 = response.RefreshToken
	})

	t.Run("refresh token 1", func(t *testing.T) {
		request := RefreshRequest{
			RefreshToken: refreshToken1,
		}

		response, statusCode, err := Refresh(request)
		if err != nil {
			t.Fatal(err)
		}

		require.NotNil(t, response)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("refresh token 1 again", func(t *testing.T) {
		request := RefreshRequest{
			RefreshToken: refreshToken1,
		}

		response, statusCode, err := Refresh(request)
		if err != nil {
			t.Fatal(err)
		}

		require.Nil(t, response)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
	})
}

func truncateAllTables() error {
	query := `
		DO $$ DECLARE
			table_name TEXT;
		BEGIN
			FOR table_name IN 
				SELECT tablename 
				FROM pg_tables 
				WHERE schemaname = 'public'
			LOOP
				EXECUTE format('TRUNCATE TABLE %I CASCADE', table_name);
			END LOOP;
		END $$;
	`

	_, err := db.Exec(context.Background(), query)

	return err
}
