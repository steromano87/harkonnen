package shooter_test

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"harkonnen/shooter"
	"os"
	"testing"
)

func TestNewContext(t *testing.T) {
	testContext := shooter.NewContext(context.Background(), zerolog.New(os.Stdout).With().Timestamp().Logger())
	assert.IsType(t, shooter.Context{}, testContext)
}
