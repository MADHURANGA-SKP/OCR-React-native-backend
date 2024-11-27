package main

import (
	"context"
	"errors"
	"net/http"
	"ocr/api"
	db "ocr/db/sqlc"
	util "ocr/util"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"golang.org/x/sync/errgroup"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
)

// main will listen to for graceful shutdown signals
var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGINT,
	syscall.SIGTERM,
}

func main() {
	//load configuration details
	config, err := util.Loadconfig(".")
	if err != nil {
		//logs a fatal error and exits
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	// Set up a context that listens for cancellation signals (Ctrl+C)
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	// Create a database connection pool
	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to the db")
	}

	// Run DB migrations
	rundDBMigration(config.MigrationURL, config.DBSource)

	// Create a new store instance
	store := db.NewStore(connPool)

	// Set up a waitgroup to manage goroutines
	waitgroup, ctx := errgroup.WithContext(ctx)
	defer stop()

	//initializes and Gin  server, start each routes to be avalible, and make it linked to database and task distribution system
	//and sepearate goroutine hellps to run this funtion with others cuncurrently
	runGinServer(ctx, waitgroup, config, store)

	// Wait for all goroutines to finish, including shutdown tasks
	if err := waitgroup.Wait(); err != nil {
		log.Fatal().Err(err).Msg("error from waitgroup")
	}
}

// run db migration when server start
func rundDBMigration(
	MigrationURL string,
	DBSource string,
) {
	//create new migrate instance
	migration, err := migrate.New(MigrationURL, DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance ")
	}
	//run migrate up and executes the database migration process
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		//If applying migrations fails logs a fatal error and terminates.
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Printf("db migrated succesfully")
}

// run runGinServer when server start
func runGinServer(
	ctx context.Context,
	waitgroup *errgroup.Group,
	config util.Config,
	store db.Store,
) {
	//create a new server instance if any error occur return the error
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	//Starts the server on the specific address if any error occur return the error
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}

	//uses to log details of each request
	httpServer := &http.Server{
		Addr:              config.HTTPServerAddress,
		ReadHeaderTimeout: 10 * time.Minute,
	}

	//Starts a background process and Waits until all goroutines  to finish
	waitgroup.Go(func() error {
		log.Info().Msgf("start http server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("http server failed to serve")
			return err
		}
		return nil
	})

	//Waits until all goroutines to finish and gracefully shut down the task processor
	waitgroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("gracefully shutdown http server ")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown http server")
			return err
		}
		log.Info().Msg("http server is stoped")
		return nil
	})
}
