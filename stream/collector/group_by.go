package collector

// NewGroupByCollector returns a collector that accumulates
// the input elements into a map whose keys are
// the result of applying the provided mapping function to the input elements,
// and values are the input elements.
func NewGroupByCollector[T any, K comparable](
	keyMapper func(T) K,
) Collector[T, map[K][]T, map[K][]T] {
	return NewGroupByCollectorWithValueMapper[T, K, T](
		keyMapper,
		Identify[T](),
	)
}

// NewGroupByCollectorWithValueMapper returns a collector that accumulates
// the input elements into a map whose keys are
// the result of applying the provided mapping function to the input elements,
func NewGroupByCollectorWithValueMapper[T any, K comparable, V any](
	keyMapper func(T) K,
	valueMapper func(T) (v V),
) Collector[T, map[K][]V, map[K][]V] {
	return NewBaseCollector[T, map[K][]V, map[K][]V](
		func() map[K][]V {
			return make(map[K][]V)
		},
		func(container map[K][]V, item T) {
			key := keyMapper(item)
			value := valueMapper(item)
			container[key] = append(container[key], value)
		},
		func(container map[K][]V) map[K][]V {
			return container
		},
	)
}
