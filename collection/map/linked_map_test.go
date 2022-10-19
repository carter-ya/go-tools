package _map

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLinkedHashMap_Put(t *testing.T) {
	var m Map[string, int] = NewLinkedHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)
	require.Equal(t, 3, m.Size())
	require.Equal(t, []string{"a", "b", "c"}, m.Keys())
	require.Equal(t, []int{1, 2, 3}, m.Values())
}

func TestLinkedHashMap_Keys(t *testing.T) {
	var m Map[string, int] = NewLinkedHashMap[string, int]()
	for i := 0; i < 10; i++ {
		m.Put(fmt.Sprintf("a%d", i), i)
	}
	require.Equal(t, 10, m.Size())
	for _, key := range m.Keys() {
		m.Remove(key)
	}
	require.Equal(t, 0, m.Size())
}
