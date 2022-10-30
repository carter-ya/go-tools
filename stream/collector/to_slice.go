package collector

import "sync"

type SliceContainer[T any] struct {
	items []T
}

func NewToSliceCollector[T any](size int) Collector[T, *SliceContainer[T], []T] {
	return NewBaseCollector[T, *SliceContainer[T], []T](
		func() *SliceContainer[T] {
			return &SliceContainer[T]{
				items: make([]T, 0, size),
			}
		},
		func(container *SliceContainer[T], item T) {
			container.items = append(container.items, item)
		},
		func(container *SliceContainer[T]) []T {
			return container.items
		},
	)
}

func NewToSliceCollectorInParallel[T any](size int) Collector[T, *SliceContainer[T], []T] {
	mu := &sync.Mutex{}
	return NewBaseCollector[T, *SliceContainer[T], []T](
		func() *SliceContainer[T] {
			return &SliceContainer[T]{
				items: make([]T, 0, size),
			}
		},
		func(container *SliceContainer[T], item T) {
			mu.Lock()
			defer mu.Unlock()

			container.items = append(container.items, item)
		},
		func(container *SliceContainer[T]) []T {
			return container.items
		},
	)
}
