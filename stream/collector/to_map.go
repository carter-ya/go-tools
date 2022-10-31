package collector

import "fmt"

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
) Collector[T, map[K]T, map[K]T] {
	return NewToMapCollectorWithDuplicateHandler[T, K, T](
		size,
		keyMapper,
		Identify[T](),
		func(duplicateKey K, existingValue, newValue T) T {
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
) Collector[T, map[K]T, map[K]T] {
	return NewToMapCollectorWithDuplicateHandler[T, K, T](
		size,
		keyMapper,
		Identify[T](),
		func(duplicateKey K, existingValue, newValue T) T {
			return newValue
		},
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
) Collector[T, map[K]V, map[K]V] {
	return NewBaseCollector[T, map[K]V, map[K]V](
		func() map[K]V {
			return make(map[K]V, size)
		},
		func(container map[K]V, item T) {
			key := keyMapper(item)
			newValue := valueMapper(item)
			if existingValue, ok := container[key]; ok && duplicateHandler != nil {
				container[key] = duplicateHandler(key, existingValue, newValue)
			} else {
				container[key] = newValue
			}
		},
		func(container map[K]V) map[K]V {
			return container
		},
	)
}
