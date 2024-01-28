package commands

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type LoginCommand struct {
	Logger          logging.ILogger
	Storage         clientEntities.IStorage
	Proxy           clientEntities.IProxy
	CryptoProvider  clientEntities.ICryptoProvider
	login, password string
}

func (c *LoginCommand) GetName() string {
	return "login"
}

func (c *LoginCommand) GetDescription() string {
	return "log in locally or on server (in case of the first local user's log in operation)"
}

func (c *LoginCommand) Validate(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("example: login <login>")
	}
	password, err := c.readPasswordMasked()
	if err != nil {
		return fmt.Errorf("failed to read password: %v", err)
	}
	c.login, c.password = args[0], password
	return nil
}

func (c *LoginCommand) Execute() cliEntities.CommandResult {
	user, err := c.Storage.ReadUser(c.login)
	if err != nil {
		return cliEntities.CommandResult{FailureMessage: err.Error()}
	}
	if user != nil {
		err := c.CryptoProvider.VerifyPassword(user, c.password)
		if err != nil {
			msg := fmt.Errorf("password verification failed: %v", err)
			return cliEntities.CommandResult{FailureMessage: msg.Error()}
		}
		return cliEntities.CommandResult{SuccessMessage: "Logged in successfully (locally)"}
	}
	c.Logger.Warningln("User %s not found locally. Going to fetch it from server", c.login)
	userData, err := c.Proxy.LogIn(c.login, c.password)
	if err != nil {
		msg := fmt.Errorf("failed to log in: %v", err)
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	newUser := &clientEntities.User{Login: c.login}
	if err := c.CryptoProvider.HashPassword(newUser, c.password); err != nil {
		msg := fmt.Errorf("failed to store user %s data locally: %v", newUser.Login, err)
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	if err := c.Storage.SaveUser(newUser); err != nil {
		msg := fmt.Errorf("failed to store user %s data locally: %v", newUser.Login, err)
		return cliEntities.CommandResult{FailureMessage: msg.Error()}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Logged in successfully (on server)",
		SessionCookie:  userData.SessionCookie,
	}
}

func (c *LoginCommand) readPasswordMasked() (string, error) {
	c.Logger.Info("Enter password: ")
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to read user password: %v", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	reader := bufio.NewReader(os.Stdin)
	var password []rune
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return "", fmt.Errorf("failed to read user password: %v", err)
		}
		switch r {
		case '\r', '\n':
			c.Logger.Warning("\n\r")
			return string(password), nil
		case '\x7f', '\b': // Backspace key
			if len(password) > 0 {
				c.Logger.Warning("\b \b") // Move back, write space to clear, and move back again
				password = password[:len(password)-1]
			}
		default:
			c.Logger.Warning("*")
			password = append(password, r)
		}
	}
}
