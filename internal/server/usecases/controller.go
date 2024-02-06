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
	cardsRepo entities.CardsRepo
	crypto    entities.ICryptoProvider
}

func NewBaseController(
	logger logging.ILogger,
	stor entities.Storage,
	userRepo entities.UserRepo,
	credsRepo entities.CredentialsRepo,
	notesRepo entities.NotesRepo,
	filesRepo entities.FilesRepo,
	cardsRepo entities.CardsRepo,
	crypto entities.ICryptoProvider,
) *BaseController {
	return &BaseController{
		logger:    logger,
		stor:      stor,
		userRepo:  userRepo,
		credsRepo: credsRepo,
		notesRepo: notesRepo,
		filesRepo: filesRepo,
		cardsRepo: cardsRepo,
		crypto:    crypto,
	}
}

func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/user/register", c.signUp)
	r.Post("/user/login", c.signIn)

	r.Post("/secrets/credentials", c.writeCredentials)
	r.Put("/secrets/credentials", c.writeCredentials)
	r.Delete("/secrets/credentials/{rowID}", c.deleteCredentials)
	r.Get("/secrets/credentials", c.getCredentials)

	r.Post("/secrets/notes", c.writeNote)
	r.Put("/secrets/notes", c.writeNote)
	r.Delete("/secrets/notes/{noteID}", c.deleteNote)
	r.Get("/secrets/notes", c.getNotes)

	r.Post("/secrets/files", c.writeFile)
	r.Put("/secrets/files", c.writeFile)
	r.Delete("/secrets/files/{fileID}", c.deleteFile)
	r.Get("/secrets/files", c.getFiles)

	r.Post("/secrets/cards", c.writeCard)
	r.Put("/secrets/cards", c.writeCard)
	r.Delete("/secrets/cards/{cardID}", c.deleteCard)
	r.Get("/secrets/cards", c.getCards)

	r.Get("/ping", c.ping)
	return r
}
