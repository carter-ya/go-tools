package collector

import (
	"sync/atomic"
)

type CountContainer struct {
	counter int
}

type CountContainerInParallel struct {
	counter *atomic.Int64
}

// NewCountCollector returns a Collector that computes the counterimum of the elements of the stream.
//
// Note: The returned Collector is not routine-safe.
func NewCountCollector[T any]() Collector[T, *CountContainer, int] {
	return NewBaseCollector[T, *CountContainer, int](
		func() *CountContainer {
			return &CountContainer{}
		},
		func(container *CountContainer, item T) {
			container.counter++
		},
		func(container *CountContainer) int {
			return container.counter
		},
	)
}

// NewCountCollectorInParallel returns a Collector that computes the counterimum of the elements of the stream.
//
// Note: The returned Collector is routine-safe.
func NewCountCollectorInParallel[T any]() Collector[T, *CountContainerInParallel, int] {
	return NewBaseCollector[T, *CountContainerInParallel, int](
		func() *CountContainerInParallel {
			return &CountContainerInParallel{counter: &atomic.Int64{}}
		},
		func(container *CountContainerInParallel, item T) {
			container.counter.Add(1)
		},
		func(container *CountContainerInParallel) int {
			return int(container.counter.Load())
		},
	)
}
