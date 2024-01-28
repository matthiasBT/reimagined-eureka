package entities

type IProxy interface {
	LogIn(login string, password string) (*UserDataResponse, error)
	Register(login string, password string) (*UserDataResponse, error)
	// SetMasterKey(byte) error
}

type UserDataResponse struct {
	SessionCookie string
}
