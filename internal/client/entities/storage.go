package entities

type IStorage interface {
	Init() error
	Shutdown()
	ReadUser(login string) (*User, error)
	SaveUser(user *User) error
}

type User struct {
	Login        string `db:"login"`
	PasswordHash []byte `db:"pwd_hash"`
}
