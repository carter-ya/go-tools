package stream

import (
	"fmt"
)

// Identify returns the given item
func Identify[V any]() func(v V) V {
	return func(v V) V {
		return v
	}
}

// MapSupplier returns a supplier of map
func MapSupplier[K comparable, V any]() SupplierFunc {
	return func() any {
		return make(map[K]V)
	}
}

// MapSupplierWithSize returns a supplier of map with the given size
func MapSupplierWithSize[K comparable, V any](size int) SupplierFunc {
	return func() any {
		return make(map[K]V, size)
	}
}

// MapAccumulator returns an accumulator function that accumulates the given items into a map.
// If the key of the item is already in the map, the panic will be raised.
//
// keyExtractor extracts the key from the given item. You can use Identify if the key-value pair are same.
func MapAccumulator[K comparable, V any](keyExtractor func(v V) K) AccumulatorFunc {
	return MapAccumulatorWithDuplicateHandler(keyExtractor, func(duplicateKey K, v1, v2 V) {
		panic(fmt.Sprintf("duplicate key: %v, v1: %v, v2: %v", duplicateKey, v1, v2))
	})
}

// MapAccumulatorWithIgnoreDuplicate returns an accumulator function that accumulates the given items into a map.
// If the key of the item is already in the map, the old item will be replaced.
//
// keyExtractor extracts the key from the given item. You can use Identify if the key-value pair are same.
func MapAccumulatorWithIgnoreDuplicate[K comparable, V any](keyExtractor func(v V) K) AccumulatorFunc {
	return MapAccumulatorWithDuplicateHandler(keyExtractor, func(duplicateKey K, v1, v2 V) {
	})
}

// MapAccumulatorWithDuplicateHandler returns an accumulator function that accumulates the given items into a map
//
// keyExtractor extracts the key from the given item. You can use Identify if the key-value pair are same.
//
// duplicateHandler handles the duplicate key
func MapAccumulatorWithDuplicateHandler[K comparable, V any](
	keyExtractor func(v V) K,
	duplicateHandler func(duplicateKey K, v1, v2 V),
) AccumulatorFunc {
	return func(identity any, item any) any {
		m := identity.(map[K]V)
		key := keyExtractor(item.(V))
		if old, ok := m[key]; ok {
			duplicateHandler(key, old, item.(V))
		}
		m[key] = item.(V)
		return m
	}
}

// SliceSupplier returns a supplier of slice
func SliceSupplier[V any]() SupplierFunc {
	return func() any {
		return make([]V, 0)
	}
}

// SliceSupplierWithSize returns a supplier of slice with the given size
func SliceSupplierWithSize[V any](size int) SupplierFunc {
	return func() any {
		return make([]V, 0, size)
	}
}

// SliceAccumulator returns an accumulator function that accumulates the given items into a slice
func SliceAccumulator[V any]() AccumulatorFunc {
	return func(identity any, item any) any {
		return append(identity.([]V), item.(V))
	}
}

// JoiningSupplier returns a supplier of string
func JoiningSupplier[V any]() SupplierFunc {
	return func() any {
		return ""
	}
}

// JoiningAccumulator returns an accumulator function that accumulates the given items into a string
func JoiningAccumulator[V any](separator string) AccumulatorFunc {
	return func(identity any, item any) any {
		s := identity.(string)
		if len(s) > 0 {
			s += separator
		}
		s += fmt.Sprint(item)
		return s
	}
}
