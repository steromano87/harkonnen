package shooter_test

import (
	"github.com/rs/zerolog"
	"github.com/steromano87/harkonnen/shooter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type VariablePoolTestSuite struct {
	suite.Suite
	logger zerolog.Logger
}

func (suite *VariablePoolTestSuite) SetupTest() {
	suite.logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func (suite *VariablePoolTestSuite) TestSetGet() {
	vp := shooter.NewVariablePool(&suite.logger)
	vp.Set("test", "testValue")
	actualValue, err := vp.Get("test")

	assert.Equal(suite.T(), "testValue", actualValue)
	assert.NoError(suite.T(), err)
}

func (suite *VariablePoolTestSuite) TestGetNonExisting() {
	vp := shooter.NewVariablePool(&suite.logger)
	actualValue, err := vp.Get("nonExisting")

	assert.Nil(suite.T(), actualValue, "Expected nil to be returned")
	assert.IsType(suite.T(), shooter.ErrVariableNotFound{}, err)
}

func (suite *VariablePoolTestSuite) TestGetString() {
	vp := shooter.NewVariablePool(&suite.logger)
	vp.Set("test", "anotherValue")
	actualValue, err := vp.Get("test")

	assert.Equal(suite.T(), "anotherValue", actualValue)
	assert.IsType(suite.T(), "", actualValue, "It is not a string")
	assert.NoError(suite.T(), err)
}

func (suite *VariablePoolTestSuite) TestGetStringNonExisting() {
	vp := shooter.NewVariablePool(&suite.logger)
	vp.Set("test", "anotherValue")
	actualValue, err := vp.GetInt("nonExisting")

	assert.Empty(suite.T(), actualValue, "Expected empty to be returned")
	assert.IsType(suite.T(), shooter.ErrVariableNotFound{}, err)
}

func (suite *VariablePoolTestSuite) TestGetStringBadType() {
	vp := shooter.NewVariablePool(&suite.logger)
	vp.Set("test", 1)
	actualValue, err := vp.GetString("test")
	assert.Equal(suite.T(), "1", actualValue)
	assert.NoError(suite.T(), err)
}

func (suite *VariablePoolTestSuite) TestGetInt() {
	vp := shooter.NewVariablePool(&suite.logger)
	vp.Set("test", 1)
	actualValue, err := vp.GetInt("test")

	assert.Equal(suite.T(), 1, actualValue)
	assert.NoError(suite.T(), err)
}

func (suite *VariablePoolTestSuite) TestGetIntNonExisting() {
	vp := shooter.NewVariablePool(&suite.logger)
	vp.Set("test", 7)
	actualValue, err := vp.GetInt("nonExisting")

	assert.Zero(suite.T(), actualValue, "Expected zero to be returned")
	assert.IsType(suite.T(), shooter.ErrVariableNotFound{}, err)
}

func (suite *VariablePoolTestSuite) TestGetIntBadType() {
	vp := shooter.NewVariablePool(&suite.logger)
	vp.Set("test", "a")
	actualValue, err := vp.GetInt("test")

	assert.Equal(suite.T(), 0, actualValue, "Expected zero-valued int")
	assert.IsType(suite.T(), shooter.ErrBadVariableCast{}, err)
}

func (suite *VariablePoolTestSuite) TestGetBool() {
	vp := shooter.NewVariablePool(&suite.logger)
	vp.Set("test", true)
	actualValue, err := vp.GetBool("test")

	assert.Equal(suite.T(), true, actualValue)
	assert.NoError(suite.T(), err)
}

func (suite *VariablePoolTestSuite) TestGetBoolNonExisting() {
	vp := shooter.NewVariablePool(&suite.logger)
	vp.Set("test", true)
	actualValue, err := vp.GetBool("nonExisting")

	assert.False(suite.T(), actualValue, "Expected false to be returned")
	assert.IsType(suite.T(), shooter.ErrVariableNotFound{}, err)
}

func (suite *VariablePoolTestSuite) TestGetBoolBadType() {
	vp := shooter.NewVariablePool(&suite.logger)
	vp.Set("test", "booleanValue")
	actualValue, err := vp.GetBool("test")

	assert.Equal(suite.T(), false, actualValue, "Expected zero-valued bool (false)")
	assert.IsType(suite.T(), shooter.ErrBadVariableCast{}, err)
}

func (suite *VariablePoolTestSuite) TestDelete() {
	vp := shooter.NewVariablePool(&suite.logger)
	vp.Set("test", "finalValue")
	vp.Delete("test")
	actualResult, err := vp.Get("test")

	assert.Nil(suite.T(), actualResult)
	assert.IsType(suite.T(), shooter.ErrVariableNotFound{}, err)
}

func TestVariablePoolTestSuite(t *testing.T) {
	suite.Run(t, new(VariablePoolTestSuite))
}
