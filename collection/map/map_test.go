package _map

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashMap_NewHashMapFromBuiltinMap(t *testing.T) {
	lm := map[string]int{
		"a": 1,
		"b": 2,
	}
	var m HashMap[string, int] = NewHashMapFromBuiltinMap[map[string]int, string, int](lm)
	require.Equal(t, 2, m.Size())
	v, found := m.Get("a")
	require.True(t, found)
	require.Equal(t, 1, v)
}

func TestHashMap_Put(t *testing.T) {
	var m Map[string, int] = NewHashMap[string, int]()
	m.Put("a", 1)
	require.Equal(t, 1, m.Size())
	v, found := m.Get("a")
	require.True(t, found)
	require.Equal(t, 1, v)
}

func TestHashMap_PutIfAbsent(t *testing.T) {
	var m Map[string, int] = NewHashMap[string, int]()
	m.PutIfAbsent("a", 1)
	require.Equal(t, 1, m.Size())
	m.PutIfAbsent("a", 2)
	require.Equal(t, 1, m.Size())
	v, found := m.Get("a")
	require.True(t, found)
	require.Equal(t, 1, v)
}

func TestHashMap_PutAll(t *testing.T) {
	var m Map[string, int] = NewHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 1)
	var m2 Map[string, int] = NewHashMap[string, int]()
	m2.Put("c", 1)
	m2.Put("d", 1)
	m.PutAll(m2)
	require.Equal(t, 4, m.Size())
}

func TestHashMap_ComputeIfAbsent(t *testing.T) {
	var m Map[string, int] = NewHashMap[string, int]()
	m.ComputeIfAbsent("a", func(key string) int {
		return 1
	})
	require.Equal(t, 1, m.Size())
	v, found := m.Get("a")
	require.True(t, found)
	require.Equal(t, 1, v)

	m.ComputeIfAbsent("a", func(key string) int {
		return 2
	})
	require.Equal(t, 1, m.Size())
	v, found = m.Get("a")
	require.True(t, found)
	require.Equal(t, 1, v)
}

func TestHashMap_ComputeIfPresent(t *testing.T) {
	var m Map[string, int] = NewHashMap[string, int]()
	m.Put("a", 1)
	m.ComputeIfPresent("a", func(key string, oldValue int) (int, RemappingAction) {
		return 2, Replace
	})
	require.Equal(t, 1, m.Size())
	v, found := m.Get("a")
	require.True(t, found)
	require.Equal(t, 2, v)

	m.ComputeIfPresent("a", func(key string, oldValue int) (int, RemappingAction) {
		return 3, Remove
	})
	require.Equal(t, 0, m.Size())
	_, found = m.Get("a")
	require.False(t, found)

	m.Put("a", 1)
	m.ComputeIfPresent("a", func(key string, oldValue int) (int, RemappingAction) {
		return 99, Noop
	})
	require.Equal(t, 1, m.Size())
	v, found = m.Get("a")
	require.True(t, found)
	require.Equal(t, 1, v)
}

func TestHashMap_Clear(t *testing.T) {
	var m Map[string, int] = NewHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 1)
	m.Clear()
	require.Equal(t, 0, m.Size())
}
