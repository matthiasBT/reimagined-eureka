package entities

import "reimagined_eureka/internal/common"

type IProxy interface {
	LogIn(login string, password string) (*UserDataResponse, error)
	Register(login string, password string, entropy *common.Entropy) (*UserDataResponse, error)
	SetSessionCookie(cookie string)

	AddCredentials(creds *common.CredentialsReq) (int, error)
	UpdateCredentials(creds *common.CredentialsReq) error
	DeleteCredentials(rowID int) error
	ReadCredentials(startID int, batchSize int) ([]*common.CredentialsReq, error)

	AddNote(note *common.NoteReq) (int, error)
	UpdateNote(note *common.NoteReq) error
	DeleteNote(rowID int) error
	ReadNotes(startID int, batchSize int) ([]*common.NoteReq, error)

	AddFile(file *common.FileReq) (int, error)
	UpdateFile(file *common.FileReq) error
	DeleteFile(rowID int) error

	AddCard(card *common.CardReq) (int, error)
	UpdateCard(card *common.CardReq) error
	DeleteCard(rowID int) error
}

type UserDataResponse struct {
	SessionCookie string
	Entropy       *common.Entropy
}
