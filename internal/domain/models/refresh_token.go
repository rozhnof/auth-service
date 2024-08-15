package models

type RefreshToken struct {
	Token string
	Hash  []byte
}
