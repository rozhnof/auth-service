package models

type UserPassword struct {
}

func NewUserPassword(data any) UserPassword {
	return UserPassword{}
}

func UserPasswordFromString(s string) UserPassword {
	return UserPassword{}
}

func (p UserPassword) String() string {
	return ""
}

func (t UserPassword) Compare(o UserPassword) bool {
	return true
}
