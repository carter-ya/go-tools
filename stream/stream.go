package stream

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
)

type Option func(s Stream)

type Stream interface {
	Map(mapper MapFunc, opts ...Option) Stream
	FlatMap(mapper FlatMapFunc, opts ...Option) Stream
	Filter(filter FilterFunc, opts ...Option) Stream
	Concat(streams []Stream, opts ...Option) Stream
	Sort(less LessFunc, opts ...Option) Stream
	Distinct(distinct DistinctFunc, opts ...Option) Stream
	Limit(limit int64, opts ...Option) Stream

	Peek(consumer ConsumeFunc, opts ...Option) Stream

	AnyMatch(match MatchFunc, opts ...Option) bool

	AllMatch(match MatchFunc, opts ...Option) bool
	NoneMatch(match MatchFunc, opts ...Option) bool
	Count(opts ...Option) int64
	Reduce(identity any, accumulator AccumulatorFunc, opts ...Option) any
	ForEach(consumer ConsumeFunc, opts ...Option)
	ToIfaceSlice(opts ...Option) []any
	Done(opts ...Option)
}

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

// Range returns a stream of int64 from start (inclusive) to end (exclusive)
//
// startInclude indicates whether start is included in the stream
//
// endExclusive indicates whether end is excluded in the stream
func Range(startInclude, endExclusive int64, opts ...Option) Stream {
	return From(func(source chan<- any) {
		for i := startInclude; i < endExclusive; i++ {
			source <- i
		}
	}, opts...)
}

func Just(items []any, opts ...Option) Stream {
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

func Concat(first Stream, other []Stream, opts ...Option) Stream {
	return first.Concat(other, opts...)
}
