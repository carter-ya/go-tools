package stream

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
)

var _ Stream = (*concurrentStream)(nil)

type concurrentStream struct {
	source      <-chan any
	parallelism uint
}

// WithSync returns an option that sets the sync of the stream
func WithSync() func(Stream) {
	return WithParallelism(1)
}

// WithParallelism returns an option that sets the parallelism of the stream
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

	if !cs.isParallel() {
		out := make(chan any)
		go func() {
			defer close(out)
			for _, s := range concatStreams {
				s.ForEach(func(item any) {
					out <- item
				})
			}
		}()
		return cs.newStream(out)
	}

	out := make(chan any, cs.parallelism)
	go func() {
		defer close(out)

		semaphore := make(chan struct{}, cs.parallelism)
		for _, s := range concatStreams {
			semaphore <- struct{}{}
			go func(s Stream) {
				defer func() {
					<-semaphore
				}()
				s.ForEach(func(item any) {
					out <- item
				})
			}(s)
		}

		// wait for all goroutines to finish
		for i := 0; uint(i) < cs.parallelism; i++ {
			semaphore <- struct{}{}
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
	cs.applyOptions(opts...)
	if !cs.isParallel() {
		seen := make(map[any]struct{})
		return cs.doStream(func(item any, out chan<- any) {
			if _, ok := seen[item]; !ok {
				seen[item] = struct{}{}
				out <- item
			}
		})
	}

	seen := new(sync.Map)
	return cs.doStream(func(item any, out chan<- any) {
		_, loaded := seen.LoadOrStore(distinct(item), struct{}{})
		if !loaded {
			out <- item
		}
	})
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

func (cs *concurrentStream) Skip(limit int64, opts ...Option) Stream {
	cs.applyOptions(opts...)

	out := make(chan any, cs.parallelism)
	go func() {
		defer close(out)

		for item := range cs.source {
			if limit <= 0 {
				out <- item
			} else {
				limit--
			}
		}
	}()
	return cs.newStream(out)
}

func (cs *concurrentStream) TakeWhile(match MatchFunc, opts ...Option) Stream {
	cs.applyOptions(opts...)

	out := make(chan any, cs.parallelism)
	go func() {
		defer close(out)

		for item := range cs.source {
			if match(item) {
				out <- item
			} else {
				go cs.drain()
				break
			}
		}
	}()
	return cs.newStream(out)
}

func (cs *concurrentStream) DropWhile(match MatchFunc, opts ...Option) Stream {
	cs.applyOptions(opts...)

	out := make(chan any, cs.parallelism)
	go func() {
		defer close(out)

		dropping := true
		for item := range cs.source {
			if dropping && match(item) {
				continue
			}
			out <- item
			dropping = false
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

	if !cs.isParallel() {
		for item := range cs.source {
			if match(item) {
				go cs.drain()
				return true
			}
		}
		return false
	}

	semaphore := make(chan struct{}, cs.parallelism)
	var result int32 = 0
	for item := range cs.source {
		semaphore <- struct{}{}
		if atomic.LoadInt32(&result) == 0 {
			go func(item any) {
				defer func() {
					<-semaphore
				}()

				if match(item) {
					atomic.StoreInt32(&result, 1)
				}
			}(item)
		} else {
			<-semaphore
			go cs.drain()
			break
		}
	}
	// wait for all goroutines to finish
	for i := 0; uint(i) < cs.parallelism; i++ {
		semaphore <- struct{}{}
	}
	return atomic.LoadInt32(&result) == 1
}

func (cs *concurrentStream) AllMatch(match MatchFunc, opts ...Option) bool {
	cs.applyOptions(opts...)

	if !cs.isParallel() {
		for item := range cs.source {
			if !match(item) {
				go cs.drain()
				return false
			}
		}
		return true
	}

	semaphore := make(chan struct{}, cs.parallelism)
	var result int32 = 1
	for item := range cs.source {
		semaphore <- struct{}{}
		if atomic.LoadInt32(&result) == 1 {
			go func(item any) {
				defer func() {
					<-semaphore
				}()

				if !match(item) {
					atomic.StoreInt32(&result, 0)
				}
			}(item)
		} else {
			<-semaphore
			go cs.drain()
			break
		}
	}
	return atomic.LoadInt32(&result) == 1
}

func (cs *concurrentStream) NoneMatch(match MatchFunc, opts ...Option) bool {
	cs.applyOptions(opts...)

	if !cs.isParallel() {
		for item := range cs.source {
			if match(item) {
				go cs.drain()
				return false
			}
		}
		return true
	}

	semaphore := make(chan struct{}, cs.parallelism)
	var result int32 = 1
	for item := range cs.source {
		semaphore <- struct{}{}
		if atomic.LoadInt32(&result) == 1 {
			go func(item any) {
				defer func() {
					<-semaphore
				}()

				if match(item) {
					atomic.StoreInt32(&result, 0)
				}
			}(item)
		} else {
			<-semaphore
			go cs.drain()
			break
		}
	}
	return atomic.LoadInt32(&result) == 1
}

func (cs *concurrentStream) FindFirst(opts ...Option) (item any, found bool) {
	cs.applyOptions(opts...)

	for item = range cs.source {
		go cs.drain()
		return item, true
	}
	return nil, false
}

func (cs *concurrentStream) Count(opts ...Option) int64 {
	return cs.Reduce(int64(0), func(identity any, item any) any {
		cnt := identity.(int64) + 1
		if cnt < 0 {
			cnt = math.MaxInt64
		}
		return cnt
	},
		copyAndAppend[Option](WithSync(), opts...)...,
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
		copyAndAppend[Option](WithSync(), opts...)...,
	)
	return ifaces.([]any)
}

func (cs *concurrentStream) ApplyOptions(opts ...Option) Stream {
	cs.applyOptions(opts...)
	return cs
}

func (cs *concurrentStream) Collect(
	supplier func() any,
	accumulator func(container, item any),
	finisher func(container any) any,
) any {
	container := supplier()
	cs.doStreamWithTerminate(func(item any) {
		accumulator(container, item)
	})
	return finisher(container)
}

func (cs *concurrentStream) Close() {
	cs.doStreamWithTerminate(func(item any) {})
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
	if !cs.isParallel() {
		return cs.doStreamWithOptionSync(fn, terminate)
	}

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

func (cs *concurrentStream) doStreamWithOptionSync(
	fn func(item any, out chan<- any),
	terminate bool,
	opts ...Option,
) *concurrentStream {
	cs.applyOptions(opts...)
	if cs.isParallel() {
		return cs.doStreamWithOption(fn, terminate)
	}

	var out chan any
	var wg *sync.WaitGroup
	if terminate {
		wg = new(sync.WaitGroup)
		wg.Add(1)
	} else {
		out = make(chan any)
	}

	go func() {
		defer func() {
			if terminate {
				wg.Done()
			} else {
				close(out)
			}
		}()

		for item := range cs.source {
			fn(item, out)
		}
	}()

	if terminate {
		wg.Wait()
		return nil
	} else {
		return cs.newStream(out)
	}
}

func (cs *concurrentStream) isParallel() bool {
	return cs.parallelism > 1
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

func copyAndAppend[T any](item T, items ...T) []T {
	newItems := make([]T, len(items)+1)
	copy(newItems, items)
	newItems[len(items)] = item
	return newItems
}
