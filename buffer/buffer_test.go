package buffer

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBufferWriter(t *testing.T) {
	buf := NewPool().Get()

	tests := []struct {
		desc string
		f    func()
		want string
	}{
		{"AppendByte", func() { buf.AppendByte('v') }, "v"},
		{"AppendString", func() { buf.AppendString("foo") }, "foo"},
		{"AppendIntPositive", func() { buf.AppendInt(42) }, "42"},
		{"AppendIntNegative", func() { buf.AppendInt(-42) }, "-42"},
		{"AppendUint", func() { buf.AppendUint(42) }, "42"},
		{"AppendBool", func() { buf.AppendBool(true) }, "true"},
		{"AppendFloat64", func() { buf.AppendFloat(3.14, 64) }, "3.14"},
		{"AppendFloat32", func() { buf.AppendFloat(float64(float32(3.14)), 32) }, "3.14"},
		{"Write", func() { buf.Write([]byte("foo")) }, "foo"},
		{"AppendTime", func() { buf.AppendTime(time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC), time.RFC3339) }, "2000-01-02T03:04:05Z"},
		{"WriteByte", func() { buf.WriteByte('v') }, "v"},
		{"WriteString", func() { buf.WriteString("foo") }, "foo"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			buf.Reset()
			tt.f()
			assert.Equal(t, tt.want, buf.String(), "Unexpected buffer.String().")
			assert.Equal(t, tt.want, string(buf.Bytes()), "Unexpected string(buffer.Bytes())")
			assert.Equal(t, len(tt.want), buf.Len(), "Unexpected buffer length")

			assert.Equal(t, _size, buf.Cap(), "Expected buffer capacity to remain constant")
		})
	}
}

func BenchmarkBuffers(b *testing.B) {
	str := strings.Repeat("a", 1024)
	slice := make([]byte, 0, 1024)
	buf := bytes.Buffer{}
	custom := NewPool().Get()

	b.Run("ByteSlice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice = append(slice, str...)
			slice = slice[:0]
		}
	})

	b.Run("ByteBuffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf.WriteString(str)
			buf.Reset()
		}
	})

	b.Run("CustomBuffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			custom.WriteString(str)
			custom.Reset()
		}
	})

}
