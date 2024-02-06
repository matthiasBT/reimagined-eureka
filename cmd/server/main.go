package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"reimagined_eureka/internal/server/adapters"
	"reimagined_eureka/internal/server/adapters/repositories"
	"reimagined_eureka/internal/server/entities"
	"reimagined_eureka/internal/server/infra/auth"
	"reimagined_eureka/internal/server/infra/config"
	"reimagined_eureka/internal/server/infra/logging"
	"reimagined_eureka/internal/server/usecases"
)

func setupServer(logger logging.ILogger, userRepo entities.UserRepo, controller *usecases.BaseController) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logging.Middleware(logger))
	r.Use(auth.Middleware(logger, userRepo))
	r.Mount("/api", controller.Route())
	return r
}

func gracefulShutdown(srv *http.Server, done chan struct{}, logger logging.ILogger) {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quitChannel
	logger.Infof("Received signal: %v\n", sig)
	done <- struct{}{}
	time.Sleep(2 * time.Second)

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err.Error())
	}
}

func main() {
	logger := logging.SetupLogger()
	conf, err := config.Read()
	if err != nil {
		panic(err)
	}
	storage := adapters.NewPGStorage(logger, conf.DatabaseDSN)
	defer storage.Shutdown()
	userRepo := repositories.NewPGUserRepo(logger, storage)
	credsRepo := repositories.NewCredentialsRepo(logger, storage)
	notesRepo := repositories.NewNotesRepo(logger, storage)
	filesRepo := repositories.NewFilesRepo(logger, storage)
	cardsRepo := repositories.NewCardsRepo(logger, storage)
	crypto := adapters.CryptoProvider{Logger: logger}
	controller := usecases.NewBaseController(
		logger, storage, userRepo, credsRepo, notesRepo, filesRepo, cardsRepo, &crypto,
	)
	r := setupServer(logger, userRepo, controller)
	srv := http.Server{Addr: conf.ServerAddr, Handler: r}

	done := make(chan struct{}, 1)

	go func() {
		logger.Infof("Launching the server at %s\n", conf.ServerAddr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	gracefulShutdown(&srv, done, logger)
}
