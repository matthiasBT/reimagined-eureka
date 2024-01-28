package entities

type CommandResult struct {
	SuccessMessage string
	FailureMessage string
	Quit           bool
	SessionCookie  string // TODO: figure out how to refresh it (store login and password?.. what if the password changes?)
}

type Command interface {
	GetName() string
	GetDescription() string
	Validate(args ...string) error
	Execute() CommandResult
}
