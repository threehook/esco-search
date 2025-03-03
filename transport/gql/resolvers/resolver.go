// Package resolvers offers resolvers to handle GQL requests
package resolvers

import (
	"github.com/threehook/esco-search/contract"
	"github.com/threehook/esco-search/transport/gql/gen"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver defines the handlers that are used to process GQL requests.
type Resolver struct {
	App contract.App
}

func (r *Resolver) defaultConfig() gen.Config {
	c := &gen.Config{Resolvers: r}
	c.NewDefaultHandler()

	return gen.Config{Resolvers: r}
}

// NewDefaultConfig configures and returns the default resolver config.
func NewDefaultConfig(vumHandler contract.App) gen.Config {
	r := &Resolver{
		App: vumHandler,
	}

	return r.defaultConfig()
}
