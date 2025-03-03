package postgres

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, &DBTestSuite{})
}

func (d *DBTestSuite) TestSomething() {
}
