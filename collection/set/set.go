package set

import "github.com/carter-ya/go-tools/stream"

type Set[E comparable] interface {
	// Add adds an element to the set.
	Add(e E) (found bool)
	// Contains returns true if the set contains the element.
	Contains(e E) bool
	// Remove removes an element from the set.
	Remove(e E) (found bool)
	// IsEmpty returns true if the set is empty.
	IsEmpty() bool
	// Size returns the number of elements in the set.
	Size() int
	// ForEach iterates over all elements in the set.
	ForEach(consumer func(e E))
	// Clear removes all elements from the set.
	Clear()

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

func NewHashSetFromSet[E comparable](s Set[E]) HashSet[E] {
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

func (h HashSet[E]) Add(e E) (found bool) {
	_, found = h[e]
	h[e] = struct{}{}
	return
}

func (h HashSet[E]) Contains(e E) bool {
	_, found := h[e]
	return found
}

func (h HashSet[E]) Remove(e E) (found bool) {
	if _, found = h[e]; found {
		delete(h, e)
	}
	return
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

func (h HashSet[E]) Clear() {
	for e := range h {
		delete(h, e)
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
