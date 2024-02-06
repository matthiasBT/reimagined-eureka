package entities

type State interface {
	GetPrompt() string
	Execute(line string) (State, CommandResult)
}
