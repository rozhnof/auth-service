package entities

import (
	"time"

	"github.com/google/uuid"
	vobjects "github.com/rozhnof/auth-service/internal/domain/value_objects"
)

type User struct {
	id       uuid.UUID
	email    string
	password vobjects.Password

	accessToken  *vobjects.AccessToken
	refreshToken *vobjects.RefreshToken
}

func NewUser(email string, password vobjects.Password) *User {
	return &User{
		id:       uuid.New(),
		email:    email,
		password: password,
	}
}

func NewExistingUser(
	id uuid.UUID,
	email string,
	password vobjects.Password,
	refreshToken *vobjects.RefreshToken,
) *User {
	return &User{
		id:           id,
		email:        email,
		password:     password,
		refreshToken: refreshToken,
	}
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() vobjects.Password {
	return u.password
}

func (u *User) RefreshToken() *vobjects.RefreshToken {
	return u.refreshToken
}

func (u *User) AccessToken() *vobjects.AccessToken {
	return u.accessToken
}

func (u *User) UpdateTokens(accessTokenTTL time.Duration, refreshTokenTTL time.Duration, secretKey []byte) error {
	payload := vobjects.AccessTokenPayload{
		UserID: u.id,
		Email:  u.email,
	}

	at, err := vobjects.NewAccessToken(accessTokenTTL, secretKey, payload)
	if err != nil {
		return err
	}

	rt, err := vobjects.NewRefreshToken(refreshTokenTTL)
	if err != nil {
		return err
	}

	u.accessToken = &at
	u.refreshToken = &rt

	return nil
}

func (u *User) CheckPassword(password string) bool {
	return u.password.Compare(password)
}
