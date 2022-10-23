package set

import (
	"github.com/carter-ya/go-tools/collection"
	"github.com/carter-ya/go-tools/stream"
)

type Set[E comparable] interface {
	collection.Collection[E]
	// Add adds an element to the set.
	Add(e E) bool
	// AddAll adds all elements from the set to the set.
	AddAll(other collection.Collection[E]) bool
	// Remove removes an element from the set.
	Remove(e E) (found bool)
	// RemoveAll removes all elements from the set.
	RemoveAll(other collection.Collection[E]) bool
	// RemoveIf removes all elements that satisfy the predicate.
	RemoveIf(predicate func(e E) bool)
	// RetainAll retains only the elements in the set that are contained in the specified set.
	RetainAll(other collection.Collection[E])
	// Clear removes all elements from the set.
	Clear()
	// Contains returns true if the set contains the element.
	Contains(e E) bool
	// ContainsAll returns true if the set contains all elements in the specified set.
	ContainsAll(other collection.Collection[E]) bool
	// IsEmpty returns true if the set is empty.
	IsEmpty() bool
	// Size returns the number of elements in the set.
	Size() int
	// ForEach iterates over all elements in the set.
	ForEach(consumer func(e E))
	// ForEachIndexed iterates over all elements in the set.
	// The consumer function returns true to stop iterating.
	ForEachIndexed(consumer func(index int, e E) (stop bool))
	// AsSlice returns the set as a slice.
	AsSlice() []E
	// Stream returns a stream of the elements in the set.
	Stream() stream.Stream
}

type HashSet[E comparable] map[E]struct{}

func NewHashSet[E comparable]() HashSet[E] {
	return make(HashSet[E])
}

func NewHashSetWithSize[E comparable](size int) HashSet[E] {
	return make(HashSet[E], size)
}

func NewHashSetFromSlice[E comparable](slice []E) HashSet[E] {
	h := NewHashSetWithSize[E](len(slice))
	for _, e := range slice {
		h.Add(e)
	}
	return h
}

func NewHashSetFromCollection[E comparable](s collection.Collection[E]) HashSet[E] {
	h := NewHashSetWithSize[E](s.Size())
	s.ForEach(func(e E) {
		h.Add(e)
	})
	return h
}

func NewHashSetFromStream[E comparable](stream stream.Stream) HashSet[E] {
	h := NewHashSet[E]()
	stream.ForEach(func(item any) {
		h.Add(item.(E))
	})
	return h
}

func (h HashSet[E]) Add(e E) bool {
	h[e] = struct{}{}
	return true
}

func (h HashSet[E]) AddAll(other collection.Collection[E]) bool {
	other.ForEach(func(e E) {
		h[e] = struct{}{}
	})
	return true
}

func (h HashSet[E]) Remove(e E) (found bool) {
	if _, found = h[e]; found {
		delete(h, e)
	}
	return
}

func (h HashSet[E]) RemoveAll(other collection.Collection[E]) bool {
	other.ForEach(func(e E) {
		delete(h, e)
	})
	return true
}

func (h HashSet[E]) RemoveIf(predicate func(e E) bool) {
	for e := range h {
		if predicate(e) {
			delete(h, e)
		}
	}
}

func (h HashSet[E]) RetainAll(other collection.Collection[E]) {
	for e := range h {
		if !other.Contains(e) {
			delete(h, e)
		}
	}
}

func (h HashSet[E]) Clear() {
	for e := range h {
		delete(h, e)
	}
}

func (h HashSet[E]) Contains(e E) bool {
	_, found := h[e]
	return found
}

func (h HashSet[E]) ContainsAll(other collection.Collection[E]) bool {
	yes := true
	other.ForEachIndexed(func(index int, e E) (stop bool) {
		yes = h.Contains(e)
		return !yes
	})
	return yes
}

func (h HashSet[E]) IsEmpty() bool {
	return len(h) == 0
}

func (h HashSet[E]) Size() int {
	return len(h)
}

func (h HashSet[E]) ForEach(consumer func(e E)) {
	for e := range h {
		consumer(e)
	}
}

func (h HashSet[E]) ForEachIndexed(consumer func(index int, e E) (stop bool)) {
	index := 0
	for e := range h {
		if consumer(index, e) {
			break
		}
		index++
	}
}

func (h HashSet[E]) AsSlice() []E {
	s := make([]E, 0, len(h))
	for e := range h {
		s = append(s, e)
	}
	return s
}

func (h HashSet[E]) Stream() stream.Stream {
	return stream.From(func(source chan<- any) {
		for e := range h {
			source <- e
		}
	})
}
