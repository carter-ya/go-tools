package collector

import (
	"golang.org/x/exp/constraints"
	"sync"
)

type SumContainer[R constraints.Integer | constraints.Float] struct {
	accumulator R
}

// NewSumCollector returns a Collector that sums the elements of the stream.
//
// Note: The returned Collector is not routine-safe.
func NewSumCollector[T any, R constraints.Integer | constraints.Float](
	mapper func(T) R) Collector[T, *SumContainer[R], R] {
	return NewBaseCollector[T, *SumContainer[R], R](
		func() *SumContainer[R] {
			return &SumContainer[R]{accumulator: 0}
		},
		func(container *SumContainer[R], item T) {
			container.accumulator += mapper(item)
		},
		func(container *SumContainer[R]) R {
			return container.accumulator
		},
	)
}

// NewSumCollectorInParallel returns a Collector that sums the elements of the stream.
//
// Note: The returned Collector is routine-safe.
func NewSumCollectorInParallel[T any, R constraints.Integer | constraints.Float](
	mapper func(T) R) Collector[T, *SumContainer[R], R] {
	mu := &sync.Mutex{}
	return NewBaseCollector[T, *SumContainer[R], R](
		func() *SumContainer[R] {
			return &SumContainer[R]{accumulator: 0}
		},
		func(container *SumContainer[R], item T) {
			mu.Lock()
			defer mu.Unlock()

			container.accumulator += mapper(item)
		},
		func(container *SumContainer[R]) R {
			return container.accumulator
		},
	)
}
