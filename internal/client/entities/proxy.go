package entities

import "reimagined_eureka/internal/common"

type IProxy interface {
	LogIn(login string, password string) (*UserDataResponse, error)
	Register(login string, password string, entropy *common.Entropy) (*UserDataResponse, error)
	SetSessionCookie(cookie string)
	AddCredentials(creds *common.CredentialsReq) (int, error)
	AddNote(creds *common.NoteReq) (int, error)
}

type UserDataResponse struct {
	SessionCookie string
	Entropy       *common.Entropy
}
