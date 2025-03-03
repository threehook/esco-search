// Package httpserver provides the default HTTP server
// and a way to create custom http servers
package httpserver

import (
	"bytes"
	"log/slog"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
)

// Options server options.
type Options struct {
	UseLogMiddleware bool
}

// DefaultServerOptions creates the default server options
//
//nolint:golint
func DefaultServerOptions() *Options {
	return &Options{
		UseLogMiddleware: true,
	}
}

// NewGinServer create a new gin server.
func NewGinServer() *gin.Engine {
	return NewGinServerWithOpts(DefaultServerOptions())
}

// NewGinServerWithOpts create new gin server with specified options.
func NewGinServerWithOpts(opts *Options) *gin.Engine {
	if opts == nil {
		return NewGinServer()
	}
	r := gin.New()
	if opts.UseLogMiddleware {
		r.Use(loggerMiddleware)
	}
	r.GET("/v1/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return r
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	r, err := w.ResponseWriter.Write(b)
	if err != nil {
		return 0, errors.Wrap(err, "writing")
	}

	return r, nil
}

func loggerMiddleware(c *gin.Context) {
	if c.Request.RequestURI == "/v1/health" {
		return
	}
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	t := time.Now()

	c.Next()
	statusCode := c.Writer.Status()
	if statusCode >= http.StatusBadRequest {
		slog.Error("unable to complete request", "body", blw.body.String(), "status", statusCode)
	}

	if time.Since(t) > time.Second {
		slog.Warn("slow request detected", "uri", c.Request.RequestURI, "duration", time.Since(t))
	}
}
