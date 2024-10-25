package value_object

type Password string

func PasswordFromString(password string) Password {
	var p Password
	return p
}

func (t Password) Compare(other Password) bool {
	return true
}

func (t Password) String() string {
	return "password"
}
