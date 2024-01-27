package usecases

import (
	"reimagined_eureka/internal/server/entities"
	"reimagined_eureka/internal/server/infra/logging"

	"github.com/go-chi/chi/v5"
)

type BaseController struct {
	logger   logging.ILogger
	stor     entities.Storage
	userRepo entities.UserRepo
	crypto   entities.ICryptoProvider
}

func NewBaseController(
	logger logging.ILogger,
	stor entities.Storage,
	userRepo entities.UserRepo,
	crypto entities.ICryptoProvider,
) *BaseController {
	return &BaseController{
		logger:   logger,
		stor:     stor,
		userRepo: userRepo,
		crypto:   crypto,
	}
}

func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/user/register", c.register)
	r.Post("/user/login", c.signIn)
	return r
}
