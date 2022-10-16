package _map

type RemappingAction int

const (
	// Replace the old value with the new value.
	Replace RemappingAction = iota
	// Remove the exists value.
	Remove
	// Noop does nothing.
	Noop
)

type Map[K comparable, V any] interface {
	// Put adds a key-value pair to the map. If the key already exists, the old value is replaced and returned.
	Put(key K, value V) (oldValue V, oldValueFound bool)
	// PutIfAbsent adds a key-value pair to the map if the key does not exist.
	PutIfAbsent(key K, newValue V)
	// PutAll adds all key-value pairs from the other map to this map.
	PutAll(other Map[K, V])
	// ComputeIfAbsent computes the value for the given key if it does not exist.
	// If the key already exists, the existing value is returned, otherwise the value is computed and stored.
	ComputeIfAbsent(key K, mapping func(key K) V) V
	// ComputeIfPresent computes the value for the given key if it exists.
	ComputeIfPresent(key K, remapping func(key K, oldValue V) (newValue V, action RemappingAction))

	// Get returns the value for the given key. If the key does not exist, the zero value is returned.
	Get(key K) (value V, found bool)
	// GetOrDefault returns the value for the given key. If the key does not exist, the default value is returned.
	GetOrDefault(key K, defaultValue V) V

	// ContainsKey returns true if the map contains the given key.
	ContainsKey(key K) bool

	// Keys returns the keys of the map.
	Keys() []K
	// Values returns the values of the map.
	Values() []V
	// ForEach iterates over all key-value pairs in the map.
	ForEach(consumer func(key K, value V))

	// Remove removes the key-value pair with the given key from the map. If the key exists, the value is returned.
	Remove(key K) (value V, found bool)
	// RemoveIf removes all key-value pairs for which the predicate returns true.
	RemoveIf(predicate func(key K, value V) bool)
	// Clear removes all key-value pairs from the map.
	Clear()

	// IsEmpty returns true if the map is empty.
	IsEmpty() bool
	// Size returns the number of key-value pairs in the map.
	Size() int

	// AsBuiltinMap returns the map as a builtin map.
	AsBuiltinMap() map[K]V
}

type HashMap[K comparable, V any] map[K]V

func NewHashMap[K comparable, V any]() HashMap[K, V] {
	return make(HashMap[K, V])
}

func NewHashMapWithSize[K comparable, V any](size int) HashMap[K, V] {
	return make(HashMap[K, V], size)
}

// NewHashMapFromBuiltinMap creates a new HashMap from a builtin map.
func NewHashMapFromBuiltinMap[M ~map[K]V, K comparable, V any](m M) HashMap[K, V] {
	return HashMap[K, V](m)
}

func (m HashMap[K, V]) Put(key K, value V) (oldValue V, oldValueFound bool) {
	oldValue, oldValueFound = m[key]
	m[key] = value
	return
}

func (m HashMap[K, V]) PutAll(other Map[K, V]) {
	other.ForEach(func(key K, value V) {
		m[key] = value
	})
}

func (m HashMap[K, V]) PutIfAbsent(key K, newValue V) {
	if _, found := m[key]; !found {
		m[key] = newValue
	}
}

func (m HashMap[K, V]) ComputeIfAbsent(key K, mapping func(key K) V) V {
	return ComputeIfAbsent(m, key, mapping)
}

func (m HashMap[K, V]) ComputeIfPresent(key K,
	remapping func(key K, oldValue V) (newValue V, action RemappingAction),
) {
	ComputeIfPresent(m, key, remapping)
}

func (m HashMap[K, V]) Get(key K) (value V, found bool) {
	value, found = m[key]
	return
}

func (m HashMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	return GetOrDefault(m, key, defaultValue)
}

func (m HashMap[K, V]) ForEach(consumer func(key K, value V)) {
	ForEach(m, consumer)
}

func (m HashMap[K, V]) ContainsKey(key K) bool {
	_, found := m[key]
	return found
}

func (m HashMap[K, V]) Keys() []K {
	return Keys(m)
}

func (m HashMap[K, V]) Values() []V {
	return Values(m)
}

func (m HashMap[K, V]) Remove(key K) (value V, found bool) {
	if value, found = m[key]; found {
		delete(m, key)
	}
	return
}

func (m HashMap[K, V]) RemoveIf(predicate func(key K, value V) bool) {
	RemoveIf(m, predicate)
}

func (m HashMap[K, V]) Clear() {
	for key := range m {
		delete(m, key)
	}
}

func (m HashMap[K, V]) IsEmpty() bool {
	return len(m) == 0
}

func (m HashMap[K, V]) Size() int {
	return len(m)
}

func (m HashMap[K, V]) AsBuiltinMap() map[K]V {
	return m
}
