package entities

type IProxy interface {
	SignIn(login string, password string) (*UserDataResponse, error)
	SignUp(login string, password string) (*UserDataResponse, error)
	// SetMasterKey(byte) error
}

type UserDataResponse struct {
	SessionCookie string
}
