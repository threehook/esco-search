package main

import (
	"log"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/threehook/esco-search/app"
	"github.com/threehook/esco-search/config"
	"github.com/threehook/esco-search/contract"
	"github.com/threehook/esco-search/internal/repository/postgres"
	"github.com/threehook/esco-search/pkg/servemux"
	"github.com/threehook/esco-search/transport/gql"
)

func initConfig() *config.Config {
	conf, err := config.New()
	if err != nil {
		log.Fatalf("loading configuration from environment failed", contract.ErrorKey, err)
	}

	return conf
}

func initApp(config *config.Config) *app.App {
	postgres.InitDatabase(config.Postgres)

	repository, err := postgres.NewRepository(config.Postgres)
	if err != nil {
		log.Fatalf("creating a database repository: %s", err.Error())
	}
	return app.NewApp(repository)
}

//nolint:gocritic //config doesn't need to be changed and size is still small enough so no pointer needed.
func initGQLServer(conf gql.Config, muxes servemux.Muxes, app contract.App) {
	mux := muxes.ForPort(conf.Port)
	server := gql.New(app, &conf)
	server.Register(mux)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	config := initConfig()

	muxes := make(servemux.Muxes)
	app := initApp(config)

	initGQLServer(config.GQLServer, muxes, app)
	fatal("running server", muxes.ListenAndServe())

}

func fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	panic(msg)
}
