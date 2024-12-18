package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rozhnof/auth-service/internal/domain"
	vobjects "github.com/rozhnof/auth-service/internal/domain/value_objects"
)

type User struct {
	id        uuid.UUID
	email     string
	password  vobjects.Password
	confirmed bool

	accessToken   *vobjects.AccessToken
	refreshToken  *vobjects.RefreshToken
	registerToken *vobjects.RegisterToken
}

func NewUser(email string, password vobjects.Password) *User {
	return &User{
		id:            uuid.New(),
		email:         email,
		password:      password,
		confirmed:     false,
		registerToken: vobjects.NewRegisterToken(),
	}
}

func NewExistingUser(
	id uuid.UUID,
	email string,
	password vobjects.Password,
	confirmed bool,
	refreshToken *vobjects.RefreshToken,
	registerToken *vobjects.RegisterToken,
) *User {
	return &User{
		id:            id,
		email:         email,
		password:      password,
		confirmed:     confirmed,
		refreshToken:  refreshToken,
		registerToken: registerToken,
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

func (u *User) Confirmed() bool {
	return u.confirmed
}

func (u *User) Confirm() error {
	if u.RegisterToken() == nil {
		return errors.Wrap(domain.ErrInvalidRegisterToken, "register token not exists")
	}

	if !u.RegisterToken().Valid() {
		return errors.Wrap(domain.ErrInvalidRegisterToken, "register token is invalid")
	}

	u.confirmed = true

	return nil
}

func (u *User) RefreshToken() *vobjects.RefreshToken {
	return u.refreshToken
}

func (u *User) RegisterToken() *vobjects.RegisterToken {
	return u.registerToken
}

func (u *User) AccessToken() *vobjects.AccessToken {
	return u.accessToken
}

func (u *User) UpdateRegisterToken() error {
	u.registerToken = vobjects.NewRegisterToken()

	return nil
}

func (u *User) RefreshTokens(accessTokenTTL time.Duration, refreshTokenTTL time.Duration, secretKey []byte) error {
	payload := vobjects.AccessTokenPayload{
		UserID: u.id,
		Email:  u.email,
	}

	at, err := vobjects.NewAccessToken(accessTokenTTL, secretKey, payload)
	if err != nil {
		return err
	}

	u.accessToken = at
	u.refreshToken = vobjects.NewRefreshToken(refreshTokenTTL)

	return nil
}

func (u *User) CheckPassword(password string) bool {
	return u.password.Compare(password)
}
