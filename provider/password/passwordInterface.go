package password

type PasswordInterface interface {
	Hash(inputPassword string) string
	Match(inputPassword string, actualPassword string) bool
}
