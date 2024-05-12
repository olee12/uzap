package exit

import (
	"os"
)

var _exit = os.Exit

func With(code int) {
	_exit(code)
}

type StubbedExit struct {
	Exited bool
	Code   int
	prev   func(code int)
}

func Stub() *StubbedExit {
	s := &StubbedExit{prev: _exit}
	_exit = s.exit
	return s
}

func WithStub(f func()) *StubbedExit {
	s := Stub()
	defer s.Unstub()
	f()
	return s
}

func (se *StubbedExit) Unstub() {
	_exit = se.prev
}

func (se *StubbedExit) exit(code int) {
	se.Exited = true
	se.Code = code
}
