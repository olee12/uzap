package stacktrace

import (
	"runtime"

	"github.com/olee12/zap/buffer"
	"github.com/olee12/zap/internal/bufferpool"
	"github.com/olee12/zap/internal/pool"
)

var _stackPool = pool.New(func() *Stack {
	return &Stack{
		storage: make([]uintptr, 64),
	}
})

type Stack struct {
	pcs    []uintptr // program counters; always a subslice of storage
	frames *runtime.Frames

	// The size of pcs varies depending on requirements:
	// it will be one if the only the first frame was requested,
	// and otherwise it will reflect the depth of the call stack.
	//
	// storage decouples the slice we need (pcs) from the slice we pool.
	// We will always allocate a reasonably large storage, but we'll use
	// only as much of it as we need.
	storage []uintptr
}

// Depth specifies how deep of a stack trace should be captured.
type Depth int

const (
	// First captures only the first frame.
	First Depth = iota

	// Full captures the entire call stack, allocating more
	// storage for it if needed.
	Full
)

// Capture captures a stack trace of the specified depth, skipping
// the provided number of frames. skip=0 identifies the caller of
// Capture.
//
// The caller must call Free on the returned stacktrace after using it.
func Capture(skip int, depth Depth) *Stack {
	stack := _stackPool.Get()

	switch depth {
	case First:
		stack.pcs = stack.storage[:1]
	case Full:
		stack.pcs = stack.storage
	}

	numFrames := runtime.Callers(
		skip+2,
		stack.pcs,
	)

	if depth == Full {
		pcs := stack.pcs

		for numFrames == len(pcs) {
			pcs = make([]uintptr, len(pcs)*2)
			numFrames = runtime.Callers(
				skip+2,
				pcs,
			)
		}
		// Discard old storage instead of returning it to the pool.
		// This will adjust the pool size over time if stack traces are
		// consistently very deep.
		stack.storage = pcs
		stack.pcs = pcs[:numFrames]

	} else {
		stack.pcs = stack.pcs[:numFrames]
	}

	stack.frames = runtime.CallersFrames(stack.pcs)

	return stack
}

// Free releases resources associated with this stacktrace
// and returns it back to the pool.
func (st *Stack) Free() {
	st.frames = nil
	st.pcs = nil
	_stackPool.Put(st)
}

// Count reports the total number of frames in this stacktrace.
// Count DOES NOT change as Next is called.
func (st *Stack) Count() int {
	return len(st.pcs)
}

// Next returns the next frame in the stack trace,
// and a boolean indicating whether there are more after it.
func (st *Stack) Next() (_ runtime.Frame, more bool) {
	return st.frames.Next()
}

func Take(skip int) string {
	stack := Capture(skip+1, Full)
	defer stack.Free()

	buffer := bufferpool.Get()
	defer buffer.Free()

	stackfmt := NewFormatter(buffer)
	stackfmt.FormatStack(stack)

	return buffer.String()

}

// Formatter formats a stack trace into a readable string representation.
type Formatter struct {
	b        *buffer.Buffer
	nonEmpty bool // whether we've written at least one frame already
}

func NewFormatter(b *buffer.Buffer) Formatter {
	return Formatter{b: b}
}

// FormatStack formats all remaining frames in the provided stacktrace -- minus
// the final runtime.main/runtime.goexit frame.
func (sf *Formatter) FormatStack(stack *Stack) {
	// Note: On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	for frame, more := stack.Next(); more; frame, more = stack.Next() {
		sf.FormatFrame(frame)
	}
}

func (sf *Formatter) FormatFrame(frame runtime.Frame) {
	if sf.nonEmpty {
		sf.b.AppendByte('\n')
	}
	sf.nonEmpty = true
	sf.b.AppendString(frame.Function)
	sf.b.AppendByte('\n')
	sf.b.AppendByte('\t')
	sf.b.AppendString(frame.File)
	sf.b.AppendByte(':')
	sf.b.AppendInt(int64(frame.Line))
}
