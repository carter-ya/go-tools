package list

import (
	"github.com/carter-ya/go-tools/collection"
	"github.com/carter-ya/go-tools/stream"
)

// List is a list of elements.
type List[E comparable] interface {
	collection.Collection[E]
	// Add adds an element to the list.
	Add(e E) bool
	// AddTo adds an element to the list at the specified index.
	AddTo(index int, e E) bool
	// AddAll adds all elements to the list.
	AddAll(other collection.Collection[E]) bool
	// AddAllTo adds all elements to the list at the specified index.
	AddAllTo(index int, other collection.Collection[E]) bool
	// Set sets the element at the specified index.
	Set(index int, e E) (old E)
	// Remove removes the first corresponding element from the list.
	Remove(e E) bool
	// RemoveAt removes the element at the specified index.
	RemoveAt(index int) E
	// RemoveAll removes all elements from the specified list.
	RemoveAll(l collection.Collection[E]) bool
	// RemoveIf removes all elements that satisfy the predicate.
	RemoveIf(predicate func(e E) bool)
	// Clear removes all elements from the list.
	Clear()
	// RetainAll retains only the elements in the list that are contained in the specified list.
	RetainAll(l collection.Collection[E])
	// Contains returns true if the list contains the element.
	Contains(e E) bool
	// ContainsAll returns true if the list contains all elements in the specified list.
	ContainsAll(l collection.Collection[E]) bool
	// IndexOf returns the index of the first occurrence of the element in the list, or -1 if the element is not found.
	IndexOf(e E) int
	// LastIndexOf returns the index of the last occurrence of the element in the list, or -1 if the element is not found.
	LastIndexOf(e E) int
	// Get returns the element at the specified index.
	Get(index int) E
	// SubList returns a list containing the elements between the specified fromIndex, inclusive, and toIndex, exclusive.
	SubList(fromIndex, toIndex int) List[E]
	// IsEmpty returns true if the list is empty.
	IsEmpty() bool
	// Size returns the number of elements in the list.
	Size() int
	// ForEach iterates over all elements in the list.
	ForEach(consumer func(e E))
	// ForEachIndexed iterates over all elements in the list.
	// The consumer function returns true to stop iterating.
	ForEachIndexed(consumer func(index int, e E) (stop bool))
	// AsSlice returns the list as a slice.
	//
	// The returned slice is a copy of the list, so changes to the returned slice will not affect the list.
	AsSlice() []E
	// Stream returns a stream of the elements in the list.
	Stream() stream.Stream
}
