package concurrent

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type Future[T any] interface {
	Get() (t T, err error)
	GetWithTimeout(timeout time.Duration) (t T, isTimeout bool, err error)
}

type holder[T any] struct {
	t   T
	err error
}

type future[T any] struct {
	h  *holder[T]
	wg *sync.WaitGroup
}

func NewFuture[T any](fn func() (t T, err error)) Future[T] {
	f := &future[T]{
		wg: new(sync.WaitGroup),
	}
	f.wg.Add(1)
	go func() {
		defer f.wg.Done()

		t, err := fn()
		h := &holder[T]{t, err}
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&f.h)), unsafe.Pointer(h))
	}()
	return f
}

func (f *future[T]) Get() (t T, err error) {
	f.wg.Wait()
	return f.h.t, f.h.err
}

func (f *future[T]) GetWithTimeout(timeout time.Duration) (t T, isTimeout bool, err error) {
	ch := make(chan int, 1)
	go func() {
		f.wg.Wait()
		ch <- 0
	}()
	select {
	case <-ch:
		return f.h.t, false, f.h.err
	case <-time.After(timeout):
		return t, true, nil
	}
}
