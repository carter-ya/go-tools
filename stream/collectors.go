package stream

import (
	"fmt"
	"github.com/carter-ya/go-tools/stream/collector"
	"golang.org/x/exp/constraints"
	"strings"
)

// NewToMapCollector returns a collector that accumulates
// the input elements into a map whose keys and values are the result of applying
// the provided mapping functions to the input elements.
//
// size is the expected size of the map.
//
// keyMapper is a function that maps the input element to a key.
//
// Note: If the input elements contain duplicate keys, will panic.
//
// Note: This collector is not routine-safe.
func NewToMapCollector[T any, K comparable](
	size int,
	keyMapper func(T) K,
) (
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) {
	return NewToMapCollectorWithDuplicateHandler[T, K, T](
		size, keyMapper, collector.Identify[T](),
		func(duplicateKey K, existingValue T, newValue T) T {
			panic(fmt.Sprintf("duplicate key: %v, existingValue: %v, newValue: %v", duplicateKey, existingValue, newValue))
		},
	)
}

// NewToMapWithIgnoreDuplicateCollector returns a collector that accumulates
// the input elements into a map whose keys and values are the result of applying
// the provided mapping functions to the input elements.
//
// size is the expected size of the map.
//
// keyMapper is a function that maps the input element to a key.
//
// Note: If the input elements contain duplicate keys, the last one will be kept.
//
// Note: This collector is not routine-safe.
func NewToMapWithIgnoreDuplicateCollector[T any, K comparable](
	size int,
	keyMapper func(T) K,
) (
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) {
	return NewToMapCollectorWithDuplicateHandler[T, K, T](
		size, keyMapper, collector.Identify[T](), nil,
	)
}

// NewToMapCollectorWithDuplicateHandler returns a collector that accumulates
// the input elements into a map whose keys and values are the result of applying
// the provided mapping functions to the input elements.
//
// size is the expected size of the map.
//
// keyMapper is a function that maps the input element to a key.
//
// valueMapper is a function that maps the input element to a value.
//
// duplicateHandler is a function that handles the duplicate key. If it is nil, will ignore the duplicate key.
//
// Note: This collector is not routine-safe.
func NewToMapCollectorWithDuplicateHandler[T any, K comparable, V any](
	size int,
	keyMapper func(T) K,
	valueMapper func(T) V,
	duplicateHandler func(duplicateKey K, existingValue V, newValue V) V,
) (
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) {
	c := collector.NewToMapCollectorWithDuplicateHandler[T, K, V](
		size, keyMapper, valueMapper, duplicateHandler,
	)

	supplier = func() any {
		return c.Supplier()()
	}
	accumulator = func(container, item any) {
		c.Accumulator()(container.(map[K]V), item.(T))
	}
	finisher = func(container any) any {
		return c.Finisher()(container.(map[K]V))
	}
	return supplier, accumulator, finisher
}

func NewToSliceCollector[T any](
	size int,
) (
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) {
	c := collector.NewToSliceCollector[T](size)

	supplier = func() any {
		return c.Supplier()()
	}
	accumulator = func(container, item any) {
		c.Accumulator()(container.(*collector.SliceContainer[T]), item.(T))
	}
	finisher = func(container any) any {
		return c.Finisher()(container.(*collector.SliceContainer[T]))
	}
	return supplier, accumulator, finisher
}

func NewJoiningCollector[T any](separator string) (
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) {
	c := collector.NewJoiningCollector[T](separator)

	supplier = func() any {
		return c.Supplier()()
	}
	accumulator = func(container, item any) {
		c.Accumulator()(container.(*strings.Builder), item.(T))
	}
	finisher = func(container any) any {
		return c.Finisher()(container.(*strings.Builder))
	}
	return supplier, accumulator, finisher
}

func NewGroupByCollector[T any, K comparable](
	keyMapper func(T) K,
) (
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) {
	c := collector.NewGroupByCollector[T, K](keyMapper)

	supplier = func() any {
		return c.Supplier()()
	}
	accumulator = func(container, item any) {
		c.Accumulator()(container.(map[K][]T), item.(T))
	}
	finisher = func(container any) any {
		return c.Finisher()(container.(map[K][]T))
	}
	return supplier, accumulator, finisher
}

func NewCountCollector[T any]() (
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) {
	c := collector.NewCountCollector[T]()

	supplier = func() any {
		return c.Supplier()()
	}
	accumulator = func(container, item any) {
		c.Accumulator()(container.(*collector.CountContainer), item.(T))
	}
	finisher = func(container any) any {
		return c.Finisher()(container.(*collector.CountContainer))
	}
	return supplier, accumulator, finisher
}

func NewSumCollector[T any, R constraints.Integer | constraints.Float](
	mapper func(T) R) (
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) {
	c := collector.NewSumCollector[T, R](mapper)

	supplier = func() any {
		return c.Supplier()()
	}
	accumulator = func(container, item any) {
		c.Accumulator()(container.(*collector.SumContainer[R]), item.(T))
	}
	finisher = func(container any) any {
		return c.Finisher()(container.(*collector.SumContainer[R]))
	}
	return supplier, accumulator, finisher
}

func NewAvgCollector[T any, R constraints.Integer | constraints.Float](
	mapper func(T) R) (
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) {
	c := collector.NewAvgCollector[T, R](mapper)

	supplier = func() any {
		return c.Supplier()()
	}
	accumulator = func(container, item any) {
		c.Accumulator()(container.(*collector.AvgContainer[R]), item.(T))
	}
	finisher = func(container any) any {
		return c.Finisher()(container.(*collector.AvgContainer[R]))
	}
	return supplier, accumulator, finisher
}

func NewMaxCollector[T any, R constraints.Ordered](
	mapper func(T) R) (
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) {
	c := collector.NewMaxCollector[T, R](mapper)

	supplier = func() any {
		return c.Supplier()()
	}
	accumulator = func(container, item any) {
		c.Accumulator()(container.(*collector.MaxContainer[R]), item.(T))
	}
	finisher = func(container any) any {
		return c.Finisher()(container.(*collector.MaxContainer[R]))
	}
	return supplier, accumulator, finisher
}

func NewMinCollector[T any, R constraints.Ordered](
	mapper func(T) R) (
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) {
	c := collector.NewMinCollector[T, R](mapper)

	supplier = func() any {
		return c.Supplier()()
	}
	accumulator = func(container, item any) {
		c.Accumulator()(container.(*collector.MinContainer[R]), item.(T))
	}
	finisher = func(container any) any {
		return c.Finisher()(container.(*collector.MinContainer[R]))
	}
	return supplier, accumulator, finisher
}
