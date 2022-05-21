package password

type PasswordInterface interface {
	Hash(inputPassword string) string
	Compare(inputPassword string, actualPassword string) bool
}
