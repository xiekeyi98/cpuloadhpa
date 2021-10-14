package cpuload

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayloadPercent(t *testing.T) {

	l := NewPayloadPercent(context.Background(), 50)
	assert.Equal(t, 50, l.targetPercent)
}
