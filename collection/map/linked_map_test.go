package _map

import (
	"encoding/json"
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

func TestLinkedHashMap_String(t *testing.T) {
	var m Map[string, int] = NewHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)
	require.Equal(t, "{a: 1, b: 2, c: 3}", m.String())
}

func TestLinkedHashMap_MarshalJSON(t *testing.T) {
	var m Map[int, int] = NewLinkedHashMap[int, int]()
	for i := 0; i < 100; i++ {
		m.Put(i, i)
	}
	bz, err := json.Marshal(m)
	require.NoError(t, err)
	t.Log(string(bz))

	var unmarshalM *LinkedHashMap[int, int]
	err = json.Unmarshal(bz, &unmarshalM)
	require.NoError(t, err)
	require.Equal(t, m, unmarshalM)

	var m2 Map[string, int] = NewLinkedHashMap[string, int]()
	for i := 0; i < 100; i++ {
		m2.Put(fmt.Sprintf("%d", i), i)
	}
	bz, err = json.Marshal(m2)
	require.NoError(t, err)
	t.Log(string(bz))
}
