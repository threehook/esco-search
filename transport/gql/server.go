// Package gql runs GQL server that receives and process GQL requests
package gql

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gin-gonic/gin"
	"github.com/threehook/esco-search/contract"
	"github.com/threehook/esco-search/pkg/httpserver"
	"github.com/threehook/esco-search/transport/gql/gen"
	"github.com/threehook/esco-search/transport/gql/resolvers"
	"github.com/vektah/gqlparser/v2/ast"
)

const (
	queryCache          = 1000
	persistedQueryCache = 100
	keepAliveSeconds    = 10
	maxUploadBytes      = 200000000
)

// Config holds the configuration parameters needed to run a graphql server.
type Config struct {
	Port int `envconfig:"default=8081,GQL_SERVER_PORT"`
}

// Server implements a graphql server.
type Server struct {
	app    contract.App
	config *Config
}

// New takes an App instance and returns a new graphql server.
func New(app contract.App, config *Config) *Server {
	return &Server{
		app:    app,
		config: config,
	}
}

// Register starts the GQL server.
func (s *Server) Register(mux *http.ServeMux) {
	serverOpts := httpserver.DefaultServerOptions()
	serverOpts.UseLogMiddleware = false

	r := httpserver.NewGinServerWithOpts(serverOpts)

	ginHandler := s.gqlHandler()

	r.POST("/gql", ginHandler)
	r.GET("/gql", ginHandler)
	mux.Handle("/", r)
}

func (s *Server) gqlHandler() gin.HandlerFunc {
	gqlResolverConfig := resolvers.NewDefaultConfig(s.app)
	h := s.newGqlServer(&gqlResolverConfig)
	h.Use(extension.FixedComplexityLimit(25)) //nolint:gomnd

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) newGqlServer(gqlResolverConfig *gen.Config) *handler.Server {
	srv := handler.New(gen.NewExecutableSchema(*gqlResolverConfig))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: keepAliveSeconds * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](queryCache))

	return srv
}
