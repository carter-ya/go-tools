package _map

import "github.com/carter-ya/go-tools/stream"

// Keys returns the keys of the map.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns the values of the map.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Copy returns a copy of the map.
func Copy[M ~map[K]V, K comparable, V any](src M) M {
	c := make(M, len(src))
	for k, v := range src {
		c[k] = v
	}
	return c
}

// CopyTo copies the map to the destination map.
func CopyTo[M ~map[K]V, K comparable, V any](src M, dst M) M {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// ForEach iterates over the map and calls the consumer function for each key-value pair.
func ForEach[M ~map[K]V, K comparable, V any](m M, consumer func(k K, v V)) {
	for k, v := range m {
		consumer(k, v)
	}
}

// GetOrDefault returns the value for the given key if it exists, otherwise returns the default value.
func GetOrDefault[M ~map[K]V, K comparable, V any](m M, key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultValue
}

// PutAll adds all key-value pairs from the other map to this map.
func PutAll[M ~map[K]V, K comparable, V any](m M, other M) {
	for k, v := range other {
		m[k] = v
	}
}

// ComputeIfAbsent computes the value for the given key if it does not exist.
func ComputeIfAbsent[M ~map[K]V, K comparable, V any](m M, key K, mapping func(k K) V) {
	if _, ok := m[key]; ok {
		return
	}
	value := mapping(key)
	m[key] = value
}

// ComputeIfPresent computes the value for the given key if it exists.
func ComputeIfPresent[M ~map[K]V, K comparable, V any](
	m M, key K,
	remapping func(key K, oldValue V) (newValue V, action RemappingAction),
) {
	if oldValue, ok := m[key]; ok {
		newValue, action := remapping(key, oldValue)
		switch action {
		case Replace:
			m[key] = newValue
		case Remove:
			delete(m, key)
		}
	}
}

// RemoveIf removes all key-value pairs for which the predicate returns true.
func RemoveIf[M ~map[K]V, K comparable, V any](m M, predicate func(key K, value V) bool) {
	// See https://go.dev/doc/go1.11#performance-compiler
	for k, v := range m {
		if predicate(k, v) {
			delete(m, k)
		}
	}
}

// KeysAsStream returns the keys of the map as a stream.
func KeysAsStream[M ~map[K]V, K comparable, V any](m M) stream.Stream {
	return stream.From(func(source chan<- any) {
		for k := range m {
			source <- k
		}
	})
}

// ValuesAsStream returns the values of the map as a stream.
func ValuesAsStream[M ~map[K]V, K comparable, V any](m M) stream.Stream {
	return stream.From(func(source chan<- any) {
		for _, v := range m {
			source <- v
		}
	})
}

// MapAsStream returns the map as a stream of key-value Pair.
func MapAsStream[M ~map[K]V, K comparable, V any](m M) stream.Stream {
	return stream.From(func(source chan<- any) {
		for k, v := range m {
			source <- Pair[K, V]{
				Key:   k,
				Value: v,
			}
		}
	})
}
