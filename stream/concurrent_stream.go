package stream

import (
	"math"
	"sort"
	"sync"
)

var _ Stream = (*concurrentStream)(nil)

type concurrentStream struct {
	source      <-chan any
	parallelism uint
}

func WithSync() func(Stream) {
	return WithParallelism(1)
}

func WithParallelism(parallelism uint) func(Stream) {
	if parallelism == 0 {
		panic("parallelism must be greater than 0")
	}
	return func(s Stream) {
		if cs, ok := s.(*concurrentStream); ok {
			cs.parallelism = parallelism
		} else {
			panic("stream: WithParallelism must be used with concurrentStream")
		}
	}
}

func (cs *concurrentStream) newStream(source <-chan any) *concurrentStream {
	return &concurrentStream{source, cs.parallelism}
}

func (cs *concurrentStream) newStreamWithSlice(slice []any) *concurrentStream {
	out := make(chan any, cs.parallelism)
	go func() {
		defer close(out)
		for _, item := range slice {
			out <- item
		}
	}()
	return cs.newStream(out)
}

func (cs *concurrentStream) Map(mapper MapFunc, opts ...Option) Stream {
	return cs.doStream(func(item any, out chan<- any) {
		out <- mapper(item)
	}, opts...)
}

func (cs *concurrentStream) FlatMap(mapper FlatMapFunc, opts ...Option) Stream {
	return cs.doStream(func(item any, out chan<- any) {
		mapper(item).ForEach(func(each any) {
			out <- each
		})
	}, opts...)
}

func (cs *concurrentStream) Filter(filter FilterFunc, opts ...Option) Stream {
	return cs.doStream(func(item any, out chan<- any) {
		if filter(item) {
			out <- item
		}
	}, opts...)
}

func (cs *concurrentStream) Concat(streams []Stream, opts ...Option) Stream {
	cs.applyOptions(opts...)

	concatStreams := append([]Stream{cs}, streams...)
	out := make(chan any, cs.parallelism)
	go func() {
		for _, s := range concatStreams {
			go func(s Stream) {
				s.ForEach(func(item any) {
					out <- item
				})
			}(s)
		}
	}()
	return cs.newStream(out)
}

func (cs *concurrentStream) Sort(less LessFunc, opts ...Option) Stream {
	iface := cs.ToIfaceSlice(opts...)
	sort.Slice(iface, func(i, j int) bool {
		return less(iface[i], iface[j])
	})
	return cs.newStreamWithSlice(iface)
}

func (cs *concurrentStream) Distinct(distinct DistinctFunc, opts ...Option) Stream {
	seen := new(sync.Map)
	return cs.doStream(func(item any, out chan<- any) {
		_, loaded := seen.LoadOrStore(distinct(item), struct{}{})
		if !loaded {
			out <- item
		}
	}, opts...)
}

func (cs *concurrentStream) Limit(limit int64, opts ...Option) Stream {
	cs.applyOptions(opts...)

	out := make(chan any, cs.parallelism)
	go func() {
		defer close(out)

		for item := range cs.source {
			if limit <= 0 {
				go cs.drain()
				break
			}
			limit--
			out <- item
		}
	}()
	return cs.newStream(out)
}

func (cs *concurrentStream) Peek(consumer ConsumeFunc, opts ...Option) Stream {
	return cs.doStream(func(item any, out chan<- any) {
		consumer(item)
		out <- item
	}, opts...)
}

func (cs *concurrentStream) AnyMatch(match MatchFunc, opts ...Option) bool {
	cs.applyOptions(opts...)

	for item := range cs.source {
		if match(item) {
			go cs.drain()
			return true
		}
	}
	return false
}

func (cs *concurrentStream) AllMatch(match MatchFunc, opts ...Option) bool {
	cs.applyOptions(opts...)

	for item := range cs.source {
		if !match(item) {
			go cs.drain()
			return false
		}
	}
	return true
}

func (cs *concurrentStream) NoneMatch(match MatchFunc, opts ...Option) bool {
	cs.applyOptions(opts...)

	for item := range cs.source {
		if match(item) {
			go cs.drain()
			return false
		}
	}
	return true
}

func (cs *concurrentStream) Count(opts ...Option) int64 {
	return cs.Reduce(int64(0), func(identity any, item any) any {
		cnt := identity.(int64) + 1
		if cnt < 0 {
			cnt = math.MaxInt64
		}
		return cnt
	},
		append([]Option{WithSync()}, opts...)...,
	).(int64)
}

func (cs *concurrentStream) Reduce(identity any, accumulator AccumulatorFunc, opts ...Option) any {
	cs.doStreamWithTerminate(func(item any) {
		identity = accumulator(identity, item)
	},
		opts...,
	)
	return identity
}

func (cs *concurrentStream) ForEach(consumer ConsumeFunc, opts ...Option) {
	cs.doStreamWithTerminate(func(item any) {
		consumer(item)
	}, opts...)
}

func (cs *concurrentStream) ToIfaceSlice(opts ...Option) []any {
	ifaces := cs.Reduce(make([]any, 0), func(identity any, item any) any {
		return append(identity.([]any), item)
	},
		append([]Option{WithSync()}, opts...)...,
	)
	return ifaces.([]any)
}

func (cs *concurrentStream) Done(opts ...Option) {
	cs.doStreamWithTerminate(func(item any) {}, opts...)
}

func (cs *concurrentStream) doStream(fn func(item any, out chan<- any), opts ...Option) *concurrentStream {
	return cs.doStreamWithOption(fn, false, opts...)
}

func (cs *concurrentStream) doStreamWithTerminate(fn func(item any), opts ...Option) {
	cs.doStreamWithOption(func(item any, out chan<- any) {
		fn(item)
	}, true, opts...)
}

// doStreamWithOption is a helper function for doStream and doStreamWithTerminate
//
// fn is the function that will be executed in parallel
//
// terminate is true if the stream should be terminated after the function is executed
func (cs *concurrentStream) doStreamWithOption(
	fn func(item any, out chan<- any),
	terminate bool,
	opts ...Option,
) *concurrentStream {
	cs.applyOptions(opts...)

	var out chan any
	var wg *sync.WaitGroup
	if terminate {
		wg = new(sync.WaitGroup)
		wg.Add(1)
	} else {
		out = make(chan any, cs.parallelism)
	}

	go func() {
		defer func() {
			if terminate {
				wg.Done()
			} else {
				close(out)
			}
		}()

		semaphore := make(chan struct{}, cs.parallelism)
		for item := range cs.source {
			// acquire semaphore
			semaphore <- struct{}{}
			go func(item any) {
				// release semaphore
				defer func() { <-semaphore }()

				fn(item, out)
			}(item)
		}

		// wait for all goroutines to finish
		for i := 0; uint(i) < cs.parallelism; i++ {
			semaphore <- struct{}{}
		}
	}()

	if terminate {
		wg.Wait()
		return nil
	} else {
		return cs.newStream(out)
	}
}

func (cs *concurrentStream) applyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(cs)
	}
}

// drain the source
func (cs *concurrentStream) drain() {
	for range cs.source {
	}
}