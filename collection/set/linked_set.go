package set

import (
	_map "github.com/carter-ya/go-tools/collection/map"
	"github.com/carter-ya/go-tools/stream"
)

type LinkedHashSet[E comparable] struct {
	linkedMap *_map.LinkedHashMap[E, struct{}]
}

func NewLinkedHashSet[E comparable]() LinkedHashSet[E] {
	return LinkedHashSet[E]{
		linkedMap: _map.NewLinkedHashMap[E, struct{}](),
	}
}

func NewLinkedHashSetWithSize[E comparable](size int) LinkedHashSet[E] {
	return LinkedHashSet[E]{
		linkedMap: _map.NewLinkedHashMapWithSize[E, struct{}](size),
	}
}

func NewLinkedHashSetFromSlice[E comparable](slice []E) LinkedHashSet[E] {
	lhs := NewLinkedHashSetWithSize[E](len(slice))
	for _, e := range slice {
		lhs.Add(e)
	}
	return lhs
}

func NewLinkedHashSetFromSet[E comparable](s Set[E]) LinkedHashSet[E] {
	lhs := NewLinkedHashSetWithSize[E](s.Size())
	s.ForEach(func(e E) {
		lhs.Add(e)
	})
	return lhs
}

func NewLinkedHashSetFromStream[E comparable](stream stream.Stream) LinkedHashSet[E] {
	lhs := NewLinkedHashSet[E]()
	stream.ForEach(func(item any) {
		lhs.Add(item.(E))
	})
	return lhs
}

func (lhs *LinkedHashSet[E]) Add(e E) (found bool) {
	_, found = lhs.linkedMap.Put(e, struct{}{})
	return
}

func (lhs *LinkedHashSet[E]) Contains(e E) bool {
	return lhs.linkedMap.ContainsKey(e)
}

func (lhs *LinkedHashSet[E]) Remove(e E) (found bool) {
	_, found = lhs.linkedMap.Remove(e)
	return
}

func (lhs *LinkedHashSet[E]) IsEmpty() bool {
	return lhs.linkedMap.IsEmpty()
}

func (lhs *LinkedHashSet[E]) Size() int {
	return lhs.linkedMap.Size()
}

func (lhs *LinkedHashSet[E]) ForEach(consumer func(e E)) {
	lhs.linkedMap.ForEach(func(key E, value struct{}) {
		consumer(key)
	})
}

func (lhs *LinkedHashSet[E]) Clear() {
	lhs.linkedMap.Clear()
}

func (lhs *LinkedHashSet[E]) AsSlice() []E {
	return lhs.linkedMap.Keys()
}

func (lhs *LinkedHashSet[E]) Stream() stream.Stream {
	return stream.From(func(source chan<- any) {
		lhs.linkedMap.ForEach(func(key E, value struct{}) {
			source <- key
		})
	})
}
