package _map

import (
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
