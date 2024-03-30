package internal

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewApplication(t *testing.T) {
	log := slog.Default()
	ctx := context.TODO()
	a := NewApplication(ctx, log)

	assert.NotNil(t, a.service)
	assert.NotNil(t, a.watcher)
}
