package value_object

import "time"

type RefreshToken string

func NewRefreshToken(timeout time.Duration) RefreshToken {
	var rt RefreshToken
	return rt
}

func RefreshTokenFromString(str string) RefreshToken {
	var rt RefreshToken
	return rt
}

func (t RefreshToken) Compare(other RefreshToken) bool {
	return true
}

func (t RefreshToken) String() string {
	return "token"
}
