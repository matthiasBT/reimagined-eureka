package usecases

import (
	"github.com/go-chi/chi/v5"

	"reimagined_eureka/internal/server/entities"
	"reimagined_eureka/internal/server/infra/logging"
)

type BaseController struct {
	logger    logging.ILogger
	stor      entities.Storage
	userRepo  entities.UserRepo
	credsRepo entities.CredentialsRepo
	notesRepo entities.NotesRepo
	filesRepo entities.FilesRepo
	crypto    entities.ICryptoProvider
}

func NewBaseController(
	logger logging.ILogger,
	stor entities.Storage,
	userRepo entities.UserRepo,
	credsRepo entities.CredentialsRepo,
	notesRepo entities.NotesRepo,
	filesRepo entities.FilesRepo,
	crypto entities.ICryptoProvider,
) *BaseController {
	return &BaseController{
		logger:    logger,
		stor:      stor,
		userRepo:  userRepo,
		credsRepo: credsRepo,
		notesRepo: notesRepo,
		filesRepo: filesRepo,
		crypto:    crypto,
	}
}

func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/user/register", c.signUp)
	r.Post("/user/login", c.signIn)
	r.Post("/secrets/credentials", c.createCredentials)
	r.Post("/secrets/notes", c.createNote)
	r.Post("/secrets/files", c.createFile)
	r.Get("/ping", c.ping)
	return r
}
