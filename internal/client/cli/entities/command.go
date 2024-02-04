package entities

type CommandResult struct {
	SuccessMessage string
	FailureMessage string
	Quit           bool
	SessionCookie  string // TODO: load from db too. TODO: figure out how to refresh it (store login and password?.. what if the password changes?)
	Login          string
	MasterKey      string
}

type Command interface {
	GetName() string
	GetDescription() string
	Validate(args ...string) error
	Execute() CommandResult
}
