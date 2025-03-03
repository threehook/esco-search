package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestNewConfig() {
}

func TestConfigTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, &ConfigTestSuite{})
}
