package app

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/suite"
)

var errTest = errors.New("error")

type AppTestSuite struct {
	suite.Suite
	app *App
	//nolint:containedctx // context in struct for testing purposes.
	ctx context.Context
}

func (a *AppTestSuite) SetupTest() {

}
