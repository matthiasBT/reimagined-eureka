package global

import (
	"fmt"

	cliEntities "reimagined_eureka/internal/client/cli/entities"
	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
)

type RefreshSessionCommand struct {
	logger          logging.ILogger
	proxy           clientEntities.IProxy
	login, password string
}

func NewRefreshSessionCommand(
	logger logging.ILogger,
	proxy clientEntities.IProxy,
	login, password string,
) *RefreshSessionCommand {
	return &RefreshSessionCommand{
		logger:   logger,
		proxy:    proxy,
		login:    login,
		password: password,
	}
}

func (c *RefreshSessionCommand) GetName() string {
	return "refresh-session"
}

func (c *RefreshSessionCommand) GetDescription() string {
	return "refresh session cookie if it's not set or has expired"
}

func (c *RefreshSessionCommand) Validate(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("example: refresh-session")
	}
	return nil
}

func (c *RefreshSessionCommand) Execute() cliEntities.CommandResult {
	userData, err := c.proxy.LogIn(c.login, c.password)
	if err != nil {
		return cliEntities.CommandResult{
			FailureMessage: fmt.Errorf("failed to refresh session cookie: %v", err).Error(),
		}
	}
	return cliEntities.CommandResult{
		SuccessMessage: "Session cookie successfully refreshed!",
		SessionCookie:  userData.SessionCookie,
	}
}
