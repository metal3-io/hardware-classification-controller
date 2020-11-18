package classifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func init() {
	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))
}

func TestCheckRangeInt(t *testing.T) {
	assert.True(t, checkRangeInt(0, 0, 99))
	assert.True(t, checkRangeInt(0, 100, 99))
	assert.True(t, checkRangeInt(1, 0, 99))
	assert.False(t, checkRangeInt(100, 0, 99))
	assert.False(t, checkRangeInt(0, 9, 99))
}
