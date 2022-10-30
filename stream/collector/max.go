package collector

import (
	"golang.org/x/exp/constraints"
	"sync"
)

type MaxContainer[R constraints.Ordered] struct {
	max *R
}

// NewMaxCollector returns a Collector that computes the maximum of the elements of the stream.
//
// Note: The returned Collector is not routine-safe.
func NewMaxCollector[T any, R constraints.Ordered](
	mapper func(T) R) Collector[T, *MaxContainer[R], R] {
	return NewBaseCollector[T, *MaxContainer[R], R](
		func() *MaxContainer[R] {
			return &MaxContainer[R]{}
		},
		func(container *MaxContainer[R], item T) {
			val := mapper(item)
			if container.max == nil || val > *container.max {
				container.max = &val
			}
		},
		func(container *MaxContainer[R]) R {
			return *container.max
		},
	)
}

// NewMaxCollectorInParallel returns a Collector that computes the maximum of the elements of the stream.
//
// Note: The returned Collector is routine-safe.
func NewMaxCollectorInParallel[T any, R constraints.Ordered](
	mapper func(T) R) Collector[T, *MaxContainer[R], R] {
	mu := &sync.Mutex{}
	return NewBaseCollector[T, *MaxContainer[R], R](
		func() *MaxContainer[R] {
			return &MaxContainer[R]{}
		},
		func(container *MaxContainer[R], item T) {
			val := mapper(item)

			mu.Lock()
			defer mu.Unlock()

			if container.max == nil || val > *container.max {
				container.max = &val
			}
		},
		func(container *MaxContainer[R]) R {
			return *container.max
		},
	)
}
