package collection

import (
	"fmt"
	"math/rand"
	"sort"
)

// Partition splits the slice into chunks of the given size.
// The last chunk may be smaller than the given size.
func Partition[T any](s []T, size int) [][]T {
	return Chunk[T](s, size)
}

// Chunk splits the slice into chunks of the given size.
// The last chunk may be smaller than the given size.
func Chunk[T any](s []T, size int) [][]T {
	if size <= 0 {
		panic("invalid chunk size")
	}

	chunks := make([][]T, 0, len(s)/size+1)
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

// Map returns a new slice with the result of applying the given function
func Map[T, R any](s []T, transformer func(T) R) []R {
	return Transform[T, R](s, transformer)
}

// Transform returns a new slice with the result of applying the given function
func Transform[T, R any](s []T, transformer func(T) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = transformer(v)
	}
	return result
}

// ToMap returns a map with the given key functions.
func ToMap[T any, K comparable](s []T, keyMapper func(T) K) map[K]T {
	return ToMapWithDuplicateKeyHandler[T, K](s, keyMapper, func(duplicateKey K, existingValue T, newValue T) T {
		panic(fmt.Sprintf("duplicate key %v found", duplicateKey))
	})
}

// ToMapWithIgnoreDuplicateKey returns a map with the given key functions.
// If duplicate key found, the new value will be kept.
func ToMapWithIgnoreDuplicateKey[T any, K comparable](s []T, keyMapper func(T) K) map[K]T {
	return ToMapWithDuplicateKeyHandler[T, K](s, keyMapper, func(duplicateKey K, existingValue T, newValue T) T {
		return newValue
	})
}

// ToMapWithDuplicateKeyHandler returns a map with the given key functions.
// If duplicate key found, the duplicate key handler will be called.
func ToMapWithDuplicateKeyHandler[T any, K comparable](
	s []T,
	keyMapper func(T) K,
	duplicateKeyHandler func(duplicateKey K, existingValue T, newValue T) T,
) map[K]T {
	m := make(map[K]T, len(s))
	for _, v := range s {
		key := keyMapper(v)
		if existingValue, ok := m[key]; ok {
			m[key] = duplicateKeyHandler(key, existingValue, v)
		} else {
			m[key] = v
		}
	}
	return m
}

// GroupBy returns a map of slices that accumulates the input elements into a map whose keys are
// the result of applying the provided mapping function to the input elements,
// and values are the input elements.
func GroupBy[T any, K comparable](s []T, keyMapper func(T) K) map[K][]T {
	m := make(map[K][]T)
	for _, v := range s {
		key := keyMapper(v)
		m[key] = append(m[key], v)
	}
	return m
}

// Filter returns a new slice with all elements that satisfy the given predicate.
func Filter[T any](s []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(s))
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Intersection returns a new slice with elements that exist in both slices.
func Intersection[T comparable](s1 []T, s2 []T) []T {
	m := make(map[T]struct{}, len(s1))
	for _, v := range s1 {
		m[v] = struct{}{}
	}

	result := make([]T, 0, len(s1))
	for _, v := range s2 {
		if _, ok := m[v]; ok {
			result = append(result, v)
		}
	}
	return result
}

// Difference returns a new slice with elements that exist in the first slice but not in the second slice.
func Difference[T comparable](s1 []T, s2 []T) []T {
	m := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		m[v] = struct{}{}
	}

	result := make([]T, 0, len(s1))
	for _, v := range s1 {
		if _, ok := m[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}

// Union returns a new slice with elements that exist in either slice.
func Union[T comparable](s1 []T, s2 []T) []T {
	m := make(map[T]struct{}, len(s1)+len(s2))
	for _, v := range s1 {
		m[v] = struct{}{}
	}
	for _, v := range s2 {
		m[v] = struct{}{}
	}

	result := make([]T, 0, len(m))
	for v := range m {
		result = append(result, v)
	}
	return result
}

// Sort sorts the slice given the provided less function.
func Sort[T any](s []T, less func(i, j int) bool) {
	sort.Slice(s, less)
}

// Reverse reverses the slice in place.
func Reverse[T any](s []T) {
	half := len(s) / 2
	for i := 0; i < half; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
}

// Shuffle shuffles the slice in place.
// See https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle.
func Shuffle[T any](s []T) {
	for i := len(s) - 1; i > 0; i-- {
		idx := rand.Intn(i + 1)
		s[i], s[idx] = s[idx], s[i]
	}
}
