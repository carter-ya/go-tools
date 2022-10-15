package stream

import "github.com/carter-ya/go-tools/collection"

// KeysAsStream returns the keys of the map as a stream.
func KeysAsStream[M ~map[K]V, K comparable, V any](m M) Stream {
	return From(func(source chan<- any) {
		for k := range m {
			source <- k
		}
	})
}

// ValuesAsStream returns the values of the map as a stream.
func ValuesAsStream[M ~map[K]V, K comparable, V any](m M) Stream {
	return From(func(source chan<- any) {
		for _, v := range m {
			source <- v
		}
	})
}

// MapAsStream returns the map as a stream of key-value Pair.
func MapAsStream[M ~map[K]V, K comparable, V any](m M) Stream {
	return From(func(source chan<- any) {
		for k, v := range m {
			source <- collection.Pair[K, V]{
				Key:   k,
				Value: v,
			}
		}
	})
}
