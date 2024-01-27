package entities

type CommandResult struct {
	SuccessMessage string
	FailureMessage string
	Quit           bool
}

type Command interface {
	GetName() string
	GetDescription() string
	Validate(args ...string) error
	Execute() CommandResult
}
