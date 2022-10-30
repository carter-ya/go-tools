package collector

import (
	"golang.org/x/exp/constraints"
	"sync"
)

type AvgContainer[R constraints.Integer | constraints.Float] struct {
	accumulator R
	counter     uint64
}

// NewAvgCollector returns a Collector that computes the average of the elements of the stream.
//
// Note: The returned Collector is not routine-safe.
func NewAvgCollector[T any, R constraints.Integer | constraints.Float](
	mapper func(T) R) Collector[T, *AvgContainer[R], R] {
	return NewBaseCollector[T, *AvgContainer[R], R](
		func() *AvgContainer[R] {
			return &AvgContainer[R]{accumulator: 0, counter: 0}
		},
		func(container *AvgContainer[R], item T) {
			container.accumulator += mapper(item)
			container.counter++
		},
		func(container *AvgContainer[R]) R {
			return container.accumulator / R(container.counter)
		},
	)
}

// NewAvgCollectorInParallel returns a Collector that computes the average of the elements of the stream.
//
// Note: The returned Collector is routine-safe.
func NewAvgCollectorInParallel[T any, R constraints.Integer | constraints.Float](
	mapper func(T) R) Collector[T, *AvgContainer[R], R] {
	mu := &sync.Mutex{}
	return NewBaseCollector[T, *AvgContainer[R], R](
		func() *AvgContainer[R] {
			return &AvgContainer[R]{accumulator: 0, counter: 0}
		},
		func(container *AvgContainer[R], item T) {
			mu.Lock()
			defer mu.Unlock()

			container.accumulator += mapper(item)
			container.counter++
		},
		func(container *AvgContainer[R]) R {
			return container.accumulator / R(container.counter)
		},
	)
}
