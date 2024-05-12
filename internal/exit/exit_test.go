package exit_test

import (
	"testing"

	"github.com/olee12/zap/internal/exit"
	"github.com/stretchr/testify/assert"
)

func TestStub(t *testing.T) {
	type want struct {
		exit bool
		code int
	}

	tests := []struct {
		f    func()
		want want
	}{
		{func() { exit.With(42) }, want{exit: true, code: 42}},
		{func() {}, want{}},
	}
	for _, tt := range tests {
		s := exit.WithStub(tt.f)

		assert.Equal(t, tt.want.exit, s.Exited, "Stub captured unexpected exit value")
		assert.Equal(t, tt.want.code, s.Code, "Stub captured unexpected exit value")
	}
}
