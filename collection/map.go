package collection

// Keys returns the keys of the map.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns the values of the map.
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// ForEach iterates over the map and calls the consumer function for each key-value pair.
func ForEach[K comparable, V any](m map[K]V, consumer func(k K, v V)) {
	for k, v := range m {
		consumer(k, v)
	}
}

// GetOrDefault returns the value for the given key if it exists, otherwise returns the default value.
func GetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultValue
}

// ComputeIfAbsent computes the value for the given key if it does not exist.
func ComputeIfAbsent[K comparable, V any](m map[K]V, key K, mapping func(k K) V) V {
	if value, ok := m[key]; ok {
		return value
	}
	value := mapping(key)
	m[key] = value
	return value
}
