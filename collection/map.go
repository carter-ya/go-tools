package collection

import "github.com/carter-ya/go-tools/stream"

// Keys returns the keys of the map.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// KeysAsStream returns the keys of the map as a stream.
func KeysAsStream[M ~map[K]V, K comparable, V any](m M) stream.Stream {
	return stream.From(func(source chan<- any) {
		for k := range m {
			source <- k
		}
	})
}

// Values returns the values of the map.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
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

// ComputeIfAbsent computes the value for the given key if it does not exist.
func ComputeIfAbsent[M ~map[K]V, K comparable, V any](m M, key K, mapping func(k K) V) V {
	if value, ok := m[key]; ok {
		return value
	}
	value := mapping(key)
	m[key] = value
	return value
}
