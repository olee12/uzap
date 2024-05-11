package buffer

import (
	"strconv"
	"time"
)

const _size = 1024

type Buffer struct {
	bs   []byte
	pool Pool
}

func (b *Buffer) AppendByte(v byte) {
	b.bs = append(b.bs, v)
}

func (b *Buffer) AppendBytes(v []byte) {
	b.bs = append(b.bs, v...)
}

func (b *Buffer) AppendString(s string) {
	b.bs = append(b.bs, s...)
}

func (b *Buffer) AppendInt(i int64) {
	b.bs = strconv.AppendInt(b.bs, i, 10)
}
func (b *Buffer) AppendUint(i uint64) {
	b.bs = strconv.AppendUint(b.bs, i, 10)
}

func (b *Buffer) AppendTime(t time.Time, layout string) {
	b.bs = t.AppendFormat(b.bs, layout)
}

func (b *Buffer) AppendBool(v bool) {
	b.bs = strconv.AppendBool(b.bs, v)
}

func (b *Buffer) AppendFloat(f float64, bitSize int) {
	b.bs = strconv.AppendFloat(b.bs, f, 'f', -1, bitSize)
}

func (b *Buffer) Len() int {
	return len(b.bs)
}

func (b *Buffer) Cap() int {
	return cap(b.bs)
}

func (b *Buffer) Bytes() []byte {
	return b.bs
}

func (b *Buffer) String() string {
	return string(b.bs)
}

func (b *Buffer) Reset() {
	b.bs = b.bs[:0]
}

func (b *Buffer) Write(bs []byte) (int, error) {
	b.bs = append(b.bs, bs...)
	return len(bs), nil
}

func (b *Buffer) WriteByte(v byte) error {
	b.AppendByte(v)
	return nil
}

func (b *Buffer) WriteString(s string) (int, error) {
	b.AppendString(s)
	return len(s), nil
}

func (b *Buffer) TrimNewline() {
	if i := len(b.bs) - 1; i >= 0 {
		if b.bs[i] == '\n' {
			b.bs = b.bs[:i]
		}
	}
}

func (b *Buffer) Free() {
	b.pool.put(b)
}
