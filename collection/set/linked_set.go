package set

import (
	"github.com/carter-ya/go-tools/collection"
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

func NewLinkedHashSetFromCollection[E comparable](c collection.Collection[E]) LinkedHashSet[E] {
	lhs := NewLinkedHashSetWithSize[E](c.Size())
	lhs.AddAll(c)
	return lhs
}

func NewLinkedHashSetFromStream[E comparable](stream stream.Stream) LinkedHashSet[E] {
	lhs := NewLinkedHashSet[E]()
	stream.ForEach(func(item any) {
		lhs.Add(item.(E))
	})
	return lhs
}

func (lhs *LinkedHashSet[E]) Add(e E) bool {
	lhs.linkedMap.Put(e, struct{}{})
	return true
}

func (lhs *LinkedHashSet[E]) AddAll(other collection.Collection[E]) bool {
	other.ForEach(func(e E) {
		lhs.linkedMap.Put(e, struct{}{})
	})
	return true
}

func (lhs *LinkedHashSet[E]) Remove(e E) (found bool) {
	_, found = lhs.linkedMap.Remove(e)
	return
}

func (lhs *LinkedHashSet[E]) RemoveAll(other collection.Collection[E]) bool {
	other.ForEach(func(e E) {
		lhs.linkedMap.Remove(e)
	})
	return true
}

func (lhs *LinkedHashSet[E]) RemoveIf(predicate func(e E) bool) {
	lhs.linkedMap.ForEach(func(e E, _ struct{}) {
		if predicate(e) {
			lhs.linkedMap.Remove(e)
		}
	})
}

func (lhs *LinkedHashSet[E]) RetainAll(other collection.Collection[E]) {
	lhs.linkedMap.ForEach(func(e E, _ struct{}) {
		if !other.Contains(e) {
			lhs.linkedMap.Remove(e)
		}
	})
}

func (lhs *LinkedHashSet[E]) Clear() {
	lhs.linkedMap.Clear()
}

func (lhs *LinkedHashSet[E]) Contains(e E) bool {
	return lhs.linkedMap.ContainsKey(e)
}

func (lhs *LinkedHashSet[E]) ContainsAll(other collection.Collection[E]) bool {
	yes := true
	other.ForEachIndexed(func(_ int, e E) (stop bool) {
		yes = lhs.linkedMap.ContainsKey(e)
		return !yes
	})
	return yes
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

func (lhs *LinkedHashSet[E]) ForEachIndexed(consumer func(index int, e E) (stop bool)) {
	lhs.linkedMap.ForEachIndexed(func(index int, key E, _ struct{}) (stop bool) {
		return consumer(index, key)
	})
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

func (lhs *LinkedHashSet[E]) String() string {
	return collection.String[E](lhs)
}

func (lhs *LinkedHashSet[E]) MarshalJSON() ([]byte, error) {
	return collection.MarshalJSON[E](lhs)
}

func (lhs *LinkedHashSet[E]) UnmarshalJSON(bytes []byte) error {
	return collection.UnmarshalJSON[E](lhs, bytes)
}
