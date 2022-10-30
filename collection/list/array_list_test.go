package list

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestArrayList_Add(t *testing.T) {
	var l List[int] = NewArrayList[int]()
	for i := 0; i < 100; i++ {
		require.Equal(t, i, l.Size())
		l.Add(i)
		require.Equal(t, i, l.Get(i))
	}
	require.Equal(t, 100, l.Size())
}
func TestArrayList_AddTo(t *testing.T) {
	l := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l.AddTo(50, 100)
	require.Equal(t, 101, l.Size())
	require.Equal(t, 100, l.Get(50))
}

func TestArrayList_AddAll(t *testing.T) {
	l := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l2 := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		l2.Add(i)
	}
	l.AddAll(l2)
	require.Equal(t, 200, l.Size())
}

func TestArrayList_Remove(t *testing.T) {
	l := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l.Remove(50)
	require.Equal(t, 99, l.Size())
	require.Equal(t, 51, l.Get(50))
}

func TestArrayList_RemoveAt(t *testing.T) {
	l := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l.RemoveAt(50)
	require.Equal(t, 99, l.Size())
	require.Equal(t, 51, l.Get(50))
}

func TestArrayList_RemoveAll(t *testing.T) {
	l := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l2 := NewArrayList[int]()
	for i := 0; i < 50; i++ {
		l2.Add(i)
	}
	l.RemoveAll(l2)
	require.Equal(t, 50, l.Size())
}

func TestArrayList_RemoveIf(t *testing.T) {
	l := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l.RemoveIf(func(e int) bool {
		return e%2 == 0
	})
	require.Equal(t, 50, l.Size())
}

func TestArrayList_Clear(t *testing.T) {
	l := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l.Clear()
	require.Equal(t, 0, l.Size())
}

func TestArrayList_RetainAll(t *testing.T) {
	l := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l2 := NewArrayList[int]()
	for i := 0; i < 10; i++ {
		l2.Add(i)
	}
	l.RetainAll(l2)
	require.Equal(t, 10, l.Size())
}

func TestArrayList_Contains(t *testing.T) {
	l := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	require.True(t, l.Contains(50))
	require.False(t, l.Contains(100))
}

func TestArrayList_ContainsAll(t *testing.T) {
	l := NewArrayList[int]()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l2 := NewArrayList[int]()
	for i := 0; i < 10; i++ {
		l2.Add(i)
	}
	require.True(t, l.ContainsAll(l2))
	l2.Add(100)
	require.False(t, l.ContainsAll(l2))
}

func TestArrayList_String(t *testing.T) {
	l := NewArrayList[int]()
	l.Add(1)
	l.Add(2)
	l.Add(3)
	require.Equal(t, "[1, 2, 3]", l.String())
}

func TestArrayList_MarshalJSON(t *testing.T) {
	var l ArrayList[int]
	bz, err := l.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "[]", string(bz))

	l.Add(1)
	l.Add(2)
	l.Add(3)
	bz, err = l.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "[1,2,3]", string(bz))

	var l2 ArrayList[string]
	l2.Add("1")
	l2.Add("2")
	l2.Add("3")
	bz, err = l2.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "[\"1\",\"2\",\"3\"]", string(bz))

	var l3 = NewArrayListFromSlice[person]([]person{
		{
			Name: "alice",
			Age:  21,
		},
		{
			Name: "bob",
			Age:  22,
		},
	})

	bz, err = l3.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "[{\"name\":\"alice\",\"age\":21},{\"name\":\"bob\",\"age\":22}]", string(bz))
}

func TestArrayList_UnmarshalJSON(t *testing.T) {
	var l = NewArrayListFromSlice[person]([]person{
		{
			Name: "alice",
			Age:  21,
		},
		{
			Name: "bob",
			Age:  22,
		},
	})
	bz, err := l.MarshalJSON()
	require.NoError(t, err)
	var l2 *ArrayList[person]
	err = json.Unmarshal(bz, &l2)
	require.NoError(t, err)
	require.Equal(t, l, l2)
}

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
