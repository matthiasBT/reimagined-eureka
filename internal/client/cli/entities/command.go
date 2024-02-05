package entities

type CommandResult struct {
	SuccessMessage string
	FailureMessage string
	Quit           bool
	SessionCookie  string
	Login          string
	Password       string
	UserID         int
	MasterKey      string
}

type Command interface {
	GetName() string
	GetDescription() string
	Validate(args ...string) error
	Execute() CommandResult
}
