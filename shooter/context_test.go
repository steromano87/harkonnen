package shooter_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"harkonnen/shooter"
	"testing"
)

func TestNewContext(t *testing.T) {
	testContext := shooter.NewContext(context.Background())
	assert.IsType(t, shooter.Context{}, testContext)
}
