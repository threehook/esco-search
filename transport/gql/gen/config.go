// Package gen contains generated GQL code and configuration
package gen

import (
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
)

// NewDefaultHandler returns a GQL handler with our own error logger.
func (c *Config) NewDefaultHandler() *handler.Server {
	return c.NewDefaultHandlerWithExtraOptions()
}

//// ExtraOptions defines the extra error logger.
//type ExtraOptions struct {
//	Logger generr.ErrorLogger
//}
//
//func newDefaultExtraOptions() *ExtraOptions {
//	return &ExtraOptions{
//		Logger: &generr.LogrusErrorLogger{},
//	}
//}

// NewDefaultHandlerWithExtraOptions configurates the GQL handler with the default settings, together with our error logger.
func (c *Config) NewDefaultHandlerWithExtraOptions() *handler.Server {
	schema := NewExecutableSchema(*c)
	srv := handler.New(schema)
	//presenter := generr.NewPresenter(extraOptions.Logger)
	//srv.SetErrorPresenter(presenter.Present)
	srv.Use(extension.Introspection{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second, //nolint:gomnd
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	return srv
}
