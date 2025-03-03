// Package config contains the app config struct and functions
package config

import (
	"github.com/cockroachdb/errors"
	"github.com/threehook/esco-search/internal/repository/postgres"
	"github.com/threehook/esco-search/transport/gql"
	"github.com/vrischmann/envconfig"
)

// Config includes the data fields that the other microservice components need to set up.
type Config struct {
	GQLServer gql.Config
	Postgres  *postgres.Config
}

// New takes its 'input' from environment variables and returns everything the microservice
// needs to serve requests, including  database configuration and a vum endpoint.
func New() (*Config, error) {
	var config Config
	err := envconfig.Init(&config)
	if err != nil {
		return nil, errors.Wrap(err, "initing")
	}

	return &config, nil
}
