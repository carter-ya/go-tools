package concurrent

import (
	"testing"
	"time"
)

func TestFuture_Get(t *testing.T) {
	f := NewFuture[int](func() (int, error) {
		time.Sleep(3 * time.Second)
		return 1, nil
	})

	for i := 0; i < 10; i++ {
		go func() {
			v, err := f.Get()
			t.Log(v, err)
		}()
	}

	time.Sleep(5 * time.Second)
}

func TestFuture_GetWithTimeout(t *testing.T) {
	f := NewFuture[int](func() (int, error) {
		time.Sleep(3 * time.Second)
		return 1, nil
	})

	for i := 0; i < 10; i++ {
		go func() {
			v, timeout, err := f.GetWithTimeout(4 * time.Second)
			t.Log(v, timeout, err)
		}()
	}

	time.Sleep(5 * time.Second)
}
