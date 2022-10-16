package stream

import "golang.org/x/exp/constraints"

type (
	// source is the channel that the generator function writes to,
	// and the generator should not close the channel
	GenerateFunc func(source chan<- any)

	MapFunc         func(item any) any
	FlatMapFunc     func(item any) Stream
	FilterFunc      func(item any) bool
	MatchFunc       FilterFunc
	LessFunc        func(a, b any) bool
	DistinctFunc    func(item any) any
	AccumulatorFunc func(identity any, item any) any
	ConsumeFunc     func(item any)
	SupplierFunc    func() any

	Option func(s Stream)
)

// Stream is the interface for a stream
type Stream interface {
	// Map applies the given mapper to each item in the stream
	Map(mapper MapFunc, opts ...Option) Stream
	// FlatMap applies the given mapper to each item in the stream
	FlatMap(mapper FlatMapFunc, opts ...Option) Stream
	// Filter filters the stream by the given predicate
	Filter(filter FilterFunc, opts ...Option) Stream
	// Concat concatenates the given streams to the current stream
	Concat(streams []Stream, opts ...Option) Stream
	// Sort sorts the stream by the given less function
	Sort(less LessFunc, opts ...Option) Stream
	// Distinct removes the duplicate items in the stream
	Distinct(distinct DistinctFunc, opts ...Option) Stream
	// Skip skips the first n items in the stream
	Skip(limit int64, opts ...Option) Stream
	// Limit limits the number of items in the stream
	Limit(limit int64, opts ...Option) Stream
	// TakeWhile takes items from the stream while the given predicate is true.
	// The first item that makes the predicate false will stop the stream.
	//
	// This is a short-circuiting terminal operation.
	TakeWhile(match MatchFunc, opts ...Option) Stream
	// DropWhile drops items from the stream while the given predicate is true.
	// The first item that makes the predicate false will start the stream.
	DropWhile(match MatchFunc, opts ...Option) Stream
	// Peek applies the given consumer to each item in the stream
	Peek(consumer ConsumeFunc, opts ...Option) Stream

	// AnyMatch returns true if any item in the stream matches the given predicate, otherwise false.
	// If the stream is empty, false is returned.
	//
	// This is a short-circuiting terminal operation.
	AnyMatch(match MatchFunc, opts ...Option) bool
	// AllMatch returns true if all items in the stream match the given predicate, otherwise false.
	// If the stream is empty, true is returned.
	//
	// This is a short-circuiting terminal operation.
	AllMatch(match MatchFunc, opts ...Option) bool
	// NoneMatch returns true if no item in the stream matches the given predicate, otherwise false.
	// If the stream is empty, true is returned.
	//
	// This is a short-circuiting terminal operation.
	NoneMatch(match MatchFunc, opts ...Option) bool
	// FindFirst returns the first item in the stream.
	// If the stream is empty, nil is returned.
	//
	// If the stream has no encounter order, then any element may be returned.
	//
	// This is a short-circuiting terminal operation.
	FindFirst(opts ...Option) (item any, found bool)
	// Count returns the number of items in the stream.
	// If the count is greater than math.MaxInt64, math.MaxInt64 is returned.
	Count(opts ...Option) int64
	// Reduce reduces the stream to a single value by the given accumulator function.
	Reduce(identity any, accumulator AccumulatorFunc, opts ...Option) any
	// ForEach applies the given consumer to each item in the stream
	ForEach(consumer ConsumeFunc, opts ...Option)
	// ToIfaceSlice returns the stream as a slice of interface{}
	ToIfaceSlice(opts ...Option) []any
	// Collect collects the stream to a supplier of the given type.
	//
	// The supplier should return a new instance of the type to collect to.
	// You can use MapSupplier to create s supplier.
	//
	// The accumulator should add the item to the supplier.
	// You can use MapAccumulator to create an accumulator.
	Collect(supplier SupplierFunc, accumulator AccumulatorFunc, opts ...Option) any
	// Close closes the stream
	Close(opts ...Option)
}

// From returns a stream from the given generator function
func From(generator GenerateFunc, opts ...Option) Stream {
	source := make(chan any)
	go func() {
		defer close(source)
		generator(source)
	}()
	cs := &concurrentStream{
		source:      source,
		parallelism: 1,
	}
	cs.applyOptions(opts...)
	return cs
}

// Range returns a stream of integer from start (inclusive) to end (exclusive)
//
// startInclude indicates whether start is included in the stream
//
// endExclusive indicates whether end is excluded in the stream
func Range[T constraints.Integer](startInclude, endExclusive T, opts ...Option) Stream {
	return From(func(source chan<- any) {
		for i := startInclude; i < endExclusive; i++ {
			source <- i
		}
	}, opts...)
}

// Just returns a stream of the given items
func Just[T any](items []T, opts ...Option) Stream {
	source := make(chan any)
	go func() {
		defer close(source)

		for _, item := range items {
			source <- item
		}
	}()
	cs := &concurrentStream{
		source:      source,
		parallelism: 1,
	}
	cs.applyOptions(opts...)
	return cs
}

// Concat concatenates the given streams to a single stream
func Concat(first Stream, other []Stream, opts ...Option) Stream {
	return first.Concat(other, opts...)
}
