package collection

import "github.com/carter-ya/go-tools/stream"

type Collection[E comparable] interface {
	// Add adds an element to the collection.
	// Returns true if the element was added.
	Add(e E) bool
	// AddAll adds all elements from the collection to the collection.
	// Returns true if the collection was modified.
	AddAll(other Collection[E]) bool
	// Remove removes an element from the collection.
	Remove(e E) bool
	// RemoveAll removes all elements from the collection.
	RemoveAll(other Collection[E]) bool
	// RemoveIf removes all elements that satisfy the predicate.
	RemoveIf(predicate func(e E) bool)
	// RetainAll retains only the elements in the collection that are contained in the specified collection.
	RetainAll(other Collection[E])
	// Clear removes all elements from the collection.
	Clear()
	// Contains returns true if the collection contains the element.
	Contains(e E) bool
	// ContainsAll returns true if the collection contains all elements in the specified collection.
	ContainsAll(other Collection[E]) bool
	// IsEmpty returns true if the collection is empty.
	IsEmpty() bool
	// Size returns the number of elements in the collection.
	Size() int
	// ForEach iterates over all elements in the collection.
	ForEach(consumer func(e E))
	// ForEachIndexed iterates over all elements in the collection.
	// The consumer function returns true to stop iterating.
	ForEachIndexed(consumer func(index int, e E) (stop bool))
	// AsSlice returns the collection as a slice.
	AsSlice() []E
	// Stream returns a stream of the elements in the collection.
	Stream() stream.Stream
}
