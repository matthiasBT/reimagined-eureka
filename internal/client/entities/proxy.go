package entities

type IProxy interface {
	SignIn(string, string) error
	SignUp(string, string) error
	// SetMasterKey(byte) error
}
