package models

type RefreshToken struct {
}

func NewRefreshToken(data any) RefreshToken {
	return RefreshToken{}
}

func RefreshTokenFromString(s string) RefreshToken {
	return RefreshToken{}
}

func (t RefreshToken) String() string {
	return ""
}

func (t RefreshToken) Compare(o RefreshToken) bool {
	return true
}
