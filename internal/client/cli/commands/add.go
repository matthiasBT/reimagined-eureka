package commands

import (
	"fmt"
	"strconv"

	"reimagined_eureka/internal/client/cli/entities"
)

type AddCommand struct {
	x, y int
}

func (c *AddCommand) GetName() string {
	return "add"
}

func (c *AddCommand) GetDescription() string {
	return "add two integers"
}

func (c *AddCommand) Validate(args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("this command takes exactly 2 arguments")
	}
	x, errX := strconv.Atoi(args[0])
	y, errY := strconv.Atoi(args[1])
	if errX != nil || errY != nil {
		return fmt.Errorf("both arguments should be integers")
	}
	c.x, c.y = x, y
	return nil
}

func (c *AddCommand) Execute() entities.CommandResult {
	return entities.CommandResult{SuccessMessage: fmt.Sprintf("%d", c.x+c.y)}
}
