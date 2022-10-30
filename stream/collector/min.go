package collector

import (
	"golang.org/x/exp/constraints"
	"sync"
)

type MinContainer[R constraints.Ordered] struct {
	min *R
}

// NewMinCollector returns a Collector that computes the minimum of the elements of the stream.
//
// Note: The returned Collector is not routine-safe.
func NewMinCollector[T any, R constraints.Ordered](
	mapper func(T) R) Collector[T, *MinContainer[R], R] {
	return NewBaseCollector[T, *MinContainer[R], R](
		func() *MinContainer[R] {
			return &MinContainer[R]{}
		},
		func(container *MinContainer[R], item T) {
			val := mapper(item)
			if container.min == nil || val < *container.min {
				container.min = &val
			}
		},
		func(container *MinContainer[R]) R {
			return *container.min
		},
	)
}

// NewMinCollectorInParallel returns a Collector that computes the minimum of the elements of the stream.
//
// Note: The returned Collector is routine-safe.
func NewMinCollectorInParallel[T any, R constraints.Ordered](
	mapper func(T) R) Collector[T, *MinContainer[R], R] {
	mu := &sync.Mutex{}
	return NewBaseCollector[T, *MinContainer[R], R](
		func() *MinContainer[R] {
			return &MinContainer[R]{}
		},
		func(container *MinContainer[R], item T) {
			val := mapper(item)

			mu.Lock()
			defer mu.Unlock()

			if container.min == nil || val < *container.min {
				container.min = &val
			}
		},
		func(container *MinContainer[R]) R {
			return *container.min
		},
	)
}
