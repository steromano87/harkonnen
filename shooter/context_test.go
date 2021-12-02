package shooter_test

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"harkonnen/shooter"
	"os"
	"testing"
)

type ContextTestSuite struct {
	suite.Suite
	logger    zerolog.Logger
	shooterID string
}

func (suite *ContextTestSuite) SetupTest() {
	suite.logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	suite.shooterID = "1"
}

func (suite *ContextTestSuite) TestNewContext() {
	testContext := shooter.NewContext(context.Background(), suite.logger, suite.shooterID)
	assert.IsType(suite.T(), shooter.Context{}, testContext)
}

func TestContextTestSuite(t *testing.T) {
	suite.Run(t, new(ContextTestSuite))
}
