package list

import (
	"github.com/carter-ya/go-tools/collection"
	"github.com/carter-ya/go-tools/stream"
)

// ArrayList is a list backed by a slice.
type ArrayList[E comparable] struct {
	data []E
}

func NewArrayList[E comparable]() *ArrayList[E] {
	return &ArrayList[E]{}
}

func NewArrayListWithSize[E comparable](size int) *ArrayList[E] {
	return &ArrayList[E]{data: make([]E, 0, size)}
}

func NewArrayListFromSlice[E comparable](slice []E) *ArrayList[E] {
	al := NewArrayListWithSize[E](len(slice))
	for _, e := range slice {
		al.Add(e)
	}
	return al
}

func NewArrayListFromCollection[E comparable](c collection.Collection[E]) *ArrayList[E] {
	al := NewArrayListWithSize[E](c.Size())
	al.AddAll(c)
	return al
}

func NewArrayListFromStream[E comparable](s stream.Stream) *ArrayList[E] {
	al := NewArrayList[E]()
	s.ForEach(func(item any) {
		al.Add(item.(E))
	})
	return al
}

func (al *ArrayList[E]) Add(e E) bool {
	al.data = append(al.data, e)
	return true
}

func (al *ArrayList[E]) AddTo(index int, e E) bool {
	al.data = append(al.data[:index], append([]E{e}, al.data[index:]...)...)
	return true
}

func (al *ArrayList[E]) AddAll(other collection.Collection[E]) bool {
	al.data = append(al.data, other.AsSlice()...)
	return true
}

func (al *ArrayList[E]) AddAllTo(index int, other collection.Collection[E]) bool {
	al.data = append(al.data[:index], append(other.AsSlice(), al.data[index:]...)...)
	return true
}

func (al *ArrayList[E]) Set(index int, e E) (old E) {
	old = al.data[index]
	al.data[index] = e
	return old
}

func (al *ArrayList[E]) Remove(e E) bool {
	for i, v := range al.data {
		if v == e {
			al.RemoveAt(i)
			return true
		}
	}
	return false
}

func (al *ArrayList[E]) RemoveAt(index int) E {
	old := al.data[index]
	al.data = append(al.data[:index], al.data[index+1:]...)
	return old
}

func (al *ArrayList[E]) RemoveAll(l collection.Collection[E]) bool {
	m := make(map[E]struct{}, l.Size())
	l.ForEach(func(e E) {
		m[e] = struct{}{}
	})

	al.RemoveIf(func(e E) bool {
		_, ok := m[e]
		return ok
	})
	return true
}

func (al *ArrayList[E]) RemoveIf(predicate func(e E) bool) {
	for i := 0; i < len(al.data); {
		if predicate(al.data[i]) {
			al.RemoveAt(i)
		} else {
			i++
		}
	}
}

func (al *ArrayList[E]) Clear() {
	al.data = nil
}

func (al *ArrayList[E]) RetainAll(l collection.Collection[E]) {
	m := make(map[E]struct{}, l.Size())
	l.ForEach(func(e E) {
		m[e] = struct{}{}
	})

	al.RemoveIf(func(e E) bool {
		_, ok := m[e]
		return !ok
	})
}

func (al *ArrayList[E]) Contains(e E) bool {
	for _, v := range al.data {
		if v == e {
			return true
		}
	}
	return false
}

func (al *ArrayList[E]) ContainsAll(l collection.Collection[E]) bool {
	yes := true
	l.ForEachIndexed(func(_ int, e E) (stop bool) {
		yes = al.Contains(e)
		return !yes
	})
	return yes
}

func (al *ArrayList[E]) IndexOf(e E) int {
	for i, v := range al.data {
		if v == e {
			return i
		}
	}
	return -1
}

func (al *ArrayList[E]) LastIndexOf(e E) int {
	for i := len(al.data); i >= 0; i-- {
		if al.data[i] == e {
			return i
		}
	}
	return -1
}

func (al *ArrayList[E]) Get(index int) E {
	return al.data[index]
}

func (al *ArrayList[E]) SubList(fromIndex, toIndex int) List[E] {
	subList := al.data[fromIndex:toIndex]
	return &ArrayList[E]{data: subList}
}

func (al *ArrayList[E]) IsEmpty() bool {
	return len(al.data) == 0
}

func (al *ArrayList[E]) Size() int {
	return len(al.data)
}

func (al *ArrayList[E]) ForEach(consumer func(e E)) {
	for _, v := range al.data {
		consumer(v)
	}
}

func (al *ArrayList[E]) ForEachIndexed(consumer func(index int, e E) (stop bool)) {
	for i, v := range al.data {
		if consumer(i, v) {
			return
		}
	}
}

func (al *ArrayList[E]) AsSlice() []E {
	cp := make([]E, len(al.data))
	copy(cp, al.data)
	return cp
}

func (al *ArrayList[E]) Stream() stream.Stream {
	return stream.Just(al.data)
}

func (al *ArrayList[E]) String() string {
	return collection.String[E](al)
}
