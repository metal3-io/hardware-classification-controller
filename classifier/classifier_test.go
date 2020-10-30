package classifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckRangeInt(t *testing.T) {
	assert.True(t, checkRangeInt(0, 0, 99))
	assert.True(t, checkRangeInt(0, 100, 99))
	assert.True(t, checkRangeInt(1, 0, 99))
	assert.False(t, checkRangeInt(100, 0, 99))
	assert.False(t, checkRangeInt(0, 9, 99))
}
