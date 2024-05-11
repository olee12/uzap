package pool_test

import (
	"runtime/debug"
	"sync"
	"testing"

	"github.com/olee12/zap/internal/pool"
	"github.com/stretchr/testify/require"
)

type pooledValue[T any] struct {
	value T
}

func TestNew(t *testing.T) {
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	p := pool.New(func() *pooledValue[string] {
		return &pooledValue[string]{
			value: "new",
		}
	})

	for i := 0; i < 1_000; i++ {
		p.Put(&pooledValue[string]{
			value: t.Name(),
		})
	}

	for i := 0; i < 10; i++ {
		func() {
			x := p.Get()
			defer p.Put(x)
			require.Equal(t, t.Name(), x.value)
		}()
	}

	for i := 0; i < 1_1000; i++ {
		p.Get()
	}

	require.Equal(t, "new", p.Get().value)
}

func TestNew_Race(t *testing.T) {
	p := pool.New(func() *pooledValue[int] {
		return &pooledValue[int]{
			value: -1,
		}
	})

	var wg sync.WaitGroup
	defer wg.Wait()

	for i := 0; i < 1_000; i++ {

		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()

			x := p.Get()
			defer p.Put(x)

			if n := x.value; n >= -1 {
				x.value = i
			}
		}()
	}
}
