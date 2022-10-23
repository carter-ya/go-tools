package _map

import "container/list"

type LinkedHashMap[K comparable, V any] struct {
	hashMap HashMap[K, V]
	list    *list.List
}

func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{
		hashMap: NewHashMap[K, V](),
		list:    list.New(),
	}
}

func NewLinkedHashMapWithSize[K comparable, V any](size int) *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{
		hashMap: NewHashMapWithSize[K, V](size),
		list:    list.New(),
	}
}

func NewLinkedHashMapFromMap[K comparable, V any](m Map[K, V]) *LinkedHashMap[K, V] {
	lm := NewLinkedHashMapWithSize[K, V](m.Size())
	lm.PutAll(m)
	return lm
}

func (m *LinkedHashMap[K, V]) Put(key K, value V) (oldValue V, oldValueFound bool) {
	oldValue, oldValueFound = m.hashMap.Put(key, value)
	if !oldValueFound {
		m.list.PushBack(key)
	}
	return
}

func (m *LinkedHashMap[K, V]) PutAll(other Map[K, V]) {
	other.ForEach(func(key K, value V) {
		m.Put(key, value)
	})
}

func (m *LinkedHashMap[K, V]) PutIfAbsent(key K, newValue V) {
	if _, found := m.hashMap[key]; !found {
		m.Put(key, newValue)
	}
}

func (m *LinkedHashMap[K, V]) ComputeIfAbsent(key K, mapping func(key K) V) {
	if _, found := m.hashMap[key]; !found {
		m.Put(key, mapping(key))
	}
}

func (m *LinkedHashMap[K, V]) ComputeIfPresent(key K,
	remapping func(key K, oldValue V) (newValue V, action RemappingAction),
) {
	if oldValue, found := m.hashMap[key]; found {
		newValue, action := remapping(key, oldValue)
		switch action {
		case Replace:
			m.hashMap.Put(key, newValue)
		case Remove:
			m.Remove(key)
		}
	}
}

func (m *LinkedHashMap[K, V]) Get(key K) (value V, found bool) {
	return m.hashMap.Get(key)
}

func (m *LinkedHashMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	return m.hashMap.GetOrDefault(key, defaultValue)
}

func (m *LinkedHashMap[K, V]) ContainsKey(key K) bool {
	return m.hashMap.ContainsKey(key)
}

func (m *LinkedHashMap[K, V]) Keys() []K {
	keys := make([]K, 0, m.hashMap.Size())
	for e := m.list.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(K))
	}
	return keys
}

func (m *LinkedHashMap[K, V]) Values() []V {
	values := make([]V, 0, m.hashMap.Size())
	for e := m.list.Front(); e != nil; e = e.Next() {
		values = append(values, m.hashMap[e.Value.(K)])
	}
	return values
}

func (m *LinkedHashMap[K, V]) ForEach(consumer func(key K, value V)) {
	for e := m.list.Front(); e != nil; e = e.Next() {
		consumer(e.Value.(K), m.hashMap[e.Value.(K)])
	}
}

func (m *LinkedHashMap[K, V]) ForEachIndexed(consumer func(index int, key K, value V) (stop bool)) {
	index := 0
	for e := m.list.Front(); e != nil; e = e.Next() {
		if consumer(index, e.Value.(K), m.hashMap[e.Value.(K)]) {
			break
		}
		index++
	}
}

func (m *LinkedHashMap[K, V]) Remove(key K) (oldValue V, oldValueFound bool) {
	oldValue, oldValueFound = m.hashMap.Remove(key)
	if oldValueFound {
		for e := m.list.Front(); e != nil; e = e.Next() {
			if e.Value == key {
				m.list.Remove(e)
				break
			}
		}
	}
	return
}

func (m *LinkedHashMap[K, V]) RemoveIf(predicate func(key K, value V) bool) {
	for e := m.list.Front(); e != nil; e = e.Next() {
		key := e.Value.(K)
		if predicate(key, m.hashMap[key]) {
			m.Remove(key)
		}
	}
}

func (m *LinkedHashMap[K, V]) Clear() {
	m.hashMap.Clear()
	m.list.Init()
}

func (m *LinkedHashMap[K, V]) IsEmpty() bool {
	return m.hashMap.IsEmpty()
}

func (m *LinkedHashMap[K, V]) Size() int {
	return m.hashMap.Size()
}

func (m *LinkedHashMap[K, V]) AsBuiltinMap() map[K]V {
	return m.hashMap.AsBuiltinMap()
}
