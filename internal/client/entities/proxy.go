package entities

import "reimagined_eureka/internal/common"

type IProxy interface {
	LogIn(login string, password string) (*UserDataResponse, error)
	Register(login string, password string, entropy *common.Entropy) (*UserDataResponse, error)
	// SetMasterKey(byte) error
}

type UserDataResponse struct {
	SessionCookie string
	Entropy       *common.Entropy
}
