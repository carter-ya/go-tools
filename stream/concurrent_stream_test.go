package stream

import (
	"github.com/carter-ya/go-tools/collection"
	"github.com/stretchr/testify/require"
	"sort"
	"sync"
	"testing"
)

func TestConcurrentStream_Map(t *testing.T) {
	expectItemsFunc := func(source chan<- any) {
		for i := 0; i < 1000; i++ {
			source <- int64(i) * 2
		}
	}
	tests := []struct {
		name        string
		stream      *concurrentStream
		expectItems []any
		ordered     bool
	}{
		{
			name:        "empty stream with no parallelism",
			stream:      Just([]any{}, WithSync()).(*concurrentStream),
			expectItems: []any{},
			ordered:     true,
		},
		{
			name:        "empty stream with parallelism",
			stream:      Just([]any{}, WithParallelism(4)).(*concurrentStream),
			expectItems: []any{},
			ordered:     true,
		},
		{
			name:        "non-empty stream with no parallelism",
			stream:      Range(0, 1000, WithSync()).(*concurrentStream),
			expectItems: From(expectItemsFunc, WithSync()).ToIfaceSlice(),
			ordered:     true,
		},
		{
			name:        "non-empty stream with parallelism",
			stream:      Range(0, 1000, WithParallelism(4)).(*concurrentStream),
			expectItems: From(expectItemsFunc, WithSync()).ToIfaceSlice(),
			ordered:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualItems := test.stream.Map(func(item any) any {
				return item.(int64) * 2
			}).ToIfaceSlice()
			if !test.ordered {
				sort.Slice(actualItems, func(i, j int) bool {
					return actualItems[i].(int64) < actualItems[j].(int64)
				})
			}
			require.ElementsMatch(t, test.expectItems, actualItems)
		})
	}
}

func TestConcurrentStream_FlatMap(t *testing.T) {
	expectItemsFunc := func(source chan<- any) {
		for i := 0; i < 1000; i++ {
			source <- int64(i)
			source <- int64(i)
		}
	}
	tests := []struct {
		name        string
		stream      *concurrentStream
		expectItems []any
		ordered     bool
	}{
		{
			name:        "empty stream with no parallelism",
			stream:      Just([]any{}, WithSync()).(*concurrentStream),
			expectItems: []any{},
			ordered:     true,
		},
		{
			name:        "empty stream with parallelism",
			stream:      Just([]any{}, WithParallelism(4)).(*concurrentStream),
			expectItems: []any{},
			ordered:     true,
		},
		{
			name:        "non-empty stream with no parallelism",
			stream:      Range(0, 1000, WithSync()).(*concurrentStream),
			expectItems: From(expectItemsFunc, WithSync()).ToIfaceSlice(),
			ordered:     true,
		},
		{
			name:        "non-empty stream with parallelism",
			stream:      Range(0, 1000, WithParallelism(4)).(*concurrentStream),
			expectItems: From(expectItemsFunc, WithSync()).ToIfaceSlice(),
			ordered:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualItems := test.stream.FlatMap(func(item any) Stream {
				return Just([]any{item, item})
			}).ToIfaceSlice()
			if !test.ordered {
				sort.Slice(actualItems, func(i, j int) bool {
					return actualItems[i].(int64) < actualItems[j].(int64)
				})
			}
			require.ElementsMatch(t, test.expectItems, actualItems)
		})
	}
}

func TestConcurrentStream_Filter(t *testing.T) {
	expectItemsFunc := func(source chan<- any) {
		for i := 0; i < 1000; i++ {
			if i%2 == 0 {
				source <- int64(i)
			}
		}
	}
	tests := []struct {
		name        string
		stream      *concurrentStream
		expectItems []any
		ordered     bool
	}{
		{
			name:        "empty stream with no parallelism",
			stream:      Just([]any{}, WithSync()).(*concurrentStream),
			expectItems: []any{},
			ordered:     true,
		},
		{
			name:        "empty stream with parallelism",
			stream:      Just([]any{}, WithParallelism(4)).(*concurrentStream),
			expectItems: []any{},
			ordered:     true,
		},
		{
			name:        "non-empty stream with no parallelism",
			stream:      Range(0, 1000, WithSync()).(*concurrentStream),
			expectItems: From(expectItemsFunc, WithSync()).ToIfaceSlice(),
			ordered:     true,
		},
		{
			name:        "non-empty stream with parallelism",
			stream:      Range(0, 1000, WithParallelism(4)).(*concurrentStream),
			expectItems: From(expectItemsFunc, WithSync()).ToIfaceSlice(),
			ordered:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualItems := test.stream.Filter(func(item any) bool {
				return item.(int64)%2 == 0
			}).ToIfaceSlice()
			if !test.ordered {
				sort.Slice(actualItems, func(i, j int) bool {
					return actualItems[i].(int64) < actualItems[j].(int64)
				})
			}
			require.ElementsMatch(t, test.expectItems, actualItems)
		})
	}
}

func TestConcurrentStream_Concat(t *testing.T) {
	tests := []struct {
		name        string
		stream      *concurrentStream
		expectItems []any
		ordered     bool
	}{
		{
			name:        "non-empty stream with no parallelism",
			stream:      Range(0, 1000, WithSync()).(*concurrentStream),
			expectItems: Range(0, 4000, WithSync()).ToIfaceSlice(),
			ordered:     true,
		},
		{
			name:        "non-empty stream with parallelism",
			stream:      Range(0, 1000, WithParallelism(4)).(*concurrentStream),
			expectItems: Range(0, 4000, WithSync()).ToIfaceSlice(),
			ordered:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualItems := test.stream.Concat([]Stream{
				Range(1000, 2000),
				Range(2000, 3000),
				Range(3000, 4000),
			}).ToIfaceSlice()
			if !test.ordered {
				sort.Slice(actualItems, func(i, j int) bool {
					return actualItems[i].(int64) < actualItems[j].(int64)
				})
			}
			require.ElementsMatch(t, test.expectItems, actualItems)
		})
	}
}

func TestConcurrentStream_Sort(t *testing.T) {
	tests := []struct {
		name        string
		parallelism uint
		stream      *concurrentStream
		expectItems []any
	}{
		{
			name:        "non-empty stream with no parallelism",
			parallelism: 1,
			stream:      Range(0, 1000, WithSync()).(*concurrentStream),
			expectItems: Range(0, 1000, WithSync()).ToIfaceSlice(),
		},
		{
			name:        "non-empty stream with parallelism",
			parallelism: 4,
			stream:      Range(0, 1000, WithSync()).(*concurrentStream),
			expectItems: Range(0, 1000, WithSync()).ToIfaceSlice(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			items := test.stream.ToIfaceSlice()
			collection.Shuffle(items)

			actualItems := Just(items, WithParallelism(test.parallelism)).Sort(func(i, j any) bool {
				return i.(int64) < j.(int64)
			}).ToIfaceSlice()
			sort.Slice(actualItems, func(i, j int) bool {
				return actualItems[i].(int64) < actualItems[j].(int64)
			})

			require.Equal(t, test.expectItems, actualItems)
		})
	}
}

func TestConcurrentStream_Distinct(t *testing.T) {
	tests := []struct {
		name        string
		parallelism uint
		stream      Stream
		expectItems []any
	}{
		{
			name:        "non-empty stream with no parallelism",
			parallelism: 1,
			stream:      Concat(Range(0, 1000, WithSync()), []Stream{Range(0, 1000, WithSync())}),
			expectItems: Range(0, 1000, WithSync()).ToIfaceSlice(),
		},
		{
			name:        "non-empty stream with parallelism",
			parallelism: 4,
			stream:      Concat(Range(0, 1000, WithSync()), []Stream{Range(0, 1000, WithSync())}),
			expectItems: Range(0, 1000, WithSync()).ToIfaceSlice(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			items := test.stream.ToIfaceSlice()
			collection.Shuffle(items)

			actualItems := Just(items, WithParallelism(test.parallelism)).Distinct(func(item any) any {
				return item
			}).ToIfaceSlice()
			sort.Slice(actualItems, func(i, j int) bool {
				return actualItems[i].(int64) < actualItems[j].(int64)
			})

			require.Equal(t, test.expectItems, actualItems)
		})
	}
}

func TestConcurrentStream_Skip(t *testing.T) {
	tests := []struct {
		name        string
		parallelism uint
		stream      Stream
		expectItems []any
	}{
		{
			name:        "non-empty stream with no parallelism",
			parallelism: 1,
			stream:      Range(0, 1000, WithSync()),
			expectItems: Range(100, 1000, WithSync()).ToIfaceSlice(),
		},
		{
			name:        "non-empty stream with parallelism",
			parallelism: 4,
			stream:      Range(0, 1000, WithSync()),
			expectItems: Range(100, 1000, WithSync()).ToIfaceSlice(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualItems := test.stream.Skip(100, WithParallelism(test.parallelism)).ToIfaceSlice()
			require.Equal(t, test.expectItems, actualItems)
		})
	}
}

func TestConcurrentStream_Limit(t *testing.T) {
	tests := []struct {
		name        string
		parallelism uint
		stream      Stream
		expectItems []any
	}{
		{
			name:        "non-empty stream with no parallelism",
			parallelism: 1,
			stream:      Range(0, 1000, WithSync()),
			expectItems: Range(0, 100, WithSync()).ToIfaceSlice(),
		},
		{
			name:        "non-empty stream with parallelism",
			parallelism: 4,
			stream:      Range(0, 1000, WithSync()),
			expectItems: Range(0, 100, WithSync()).ToIfaceSlice(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualItems := test.stream.Limit(100, WithParallelism(test.parallelism)).ToIfaceSlice()
			require.Equal(t, test.expectItems, actualItems)
		})
	}
}

func TestConcurrentStream_TakeWhile(t *testing.T) {
	tests := []struct {
		name        string
		stream      Stream
		expectItems []any
		ordered     bool
	}{
		{
			name:        "non-empty stream with no parallelism",
			stream:      Range(0, 1000, WithSync()),
			expectItems: Range(0, 100, WithSync()).ToIfaceSlice(),
			ordered:     true,
		},
		{
			name:        "non-empty stream with parallelism",
			stream:      Range(0, 1000, WithParallelism(4)),
			expectItems: Range(0, 100, WithSync()).ToIfaceSlice(),
			ordered:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualItems := test.stream.TakeWhile(func(item any) bool {
				return item.(int64) < 100
			}).ToIfaceSlice()
			if !test.ordered {
				sort.Slice(actualItems, func(i, j int) bool {
					return actualItems[i].(int64) < actualItems[j].(int64)
				})
			}
			require.Equal(t, test.expectItems, actualItems)
		})
	}
}

func TestConcurrentStream_DropWhile(t *testing.T) {
	tests := []struct {
		name        string
		stream      Stream
		expectItems []any
		ordered     bool
	}{
		{
			name:        "non-empty stream with no parallelism",
			stream:      Range(0, 1000, WithSync()),
			expectItems: Range(100, 1000, WithSync()).ToIfaceSlice(),
			ordered:     true,
		},
		{
			name:        "non-empty stream with parallelism",
			stream:      Range(0, 1000, WithParallelism(4)),
			expectItems: Range(100, 1000, WithSync()).ToIfaceSlice(),
			ordered:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualItems := test.stream.DropWhile(func(item any) bool {
				return item.(int64) < 100
			}).ToIfaceSlice()
			if !test.ordered {
				sort.Slice(actualItems, func(i, j int) bool {
					return actualItems[i].(int64) < actualItems[j].(int64)
				})
			}
			require.Equal(t, test.expectItems, actualItems)
		})
	}
}

func TestConcurrentStream_Peek(t *testing.T) {
	tests := []struct {
		name        string
		stream      Stream
		expectItems []any
		ordered     bool
	}{
		{
			name:        "non-empty stream with no parallelism",
			stream:      Range(0, 1000, WithSync()),
			expectItems: Range(0, 1000, WithSync()).ToIfaceSlice(),
			ordered:     true,
		},
		{
			name:        "non-empty stream with parallelism",
			stream:      Range(0, 1000, WithParallelism(4)),
			expectItems: Range(0, 1000, WithSync()).ToIfaceSlice(),
			ordered:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			peekItems := make([]any, 0, len(test.expectItems))
			lock := new(sync.Mutex)
			actualItems := test.stream.Peek(func(item any) {
				lock.Lock()
				defer lock.Unlock()
				peekItems = append(peekItems, item)
			}).ToIfaceSlice()
			if !test.ordered {
				sort.Slice(peekItems, func(i, j int) bool {
					return peekItems[i].(int64) < peekItems[j].(int64)
				})
				sort.Slice(actualItems, func(i, j int) bool {
					return actualItems[i].(int64) < actualItems[j].(int64)
				})
			}
			require.Equal(t, test.expectItems, actualItems)
			require.Equal(t, test.expectItems, peekItems)
		})
	}
}

func TestConcurrentStream_AnyMatch(t *testing.T) {
	tests := []struct {
		name   string
		stream Stream
		expect bool
	}{
		{
			name:   "empty stream with no parallelism",
			stream: Just([]any{}, WithSync()),
			expect: false,
		},
		{
			name:   "empty stream with parallelism",
			stream: Just([]any{}, WithParallelism(4)),
			expect: false,
		},
		{
			name:   "non-empty stream with no parallelism and no match",
			stream: Range(0, 500, WithSync()),
			expect: false,
		},
		{
			name:   "non-empty stream with parallelism and no match",
			stream: Range(0, 500, WithParallelism(4)),
			expect: false,
		},
		{
			name:   "non-empty stream with no parallelism and match",
			stream: Range(0, 1e10, WithSync()),
			expect: true,
		},
		{
			name:   "non-empty stream with parallelism and match",
			stream: Range(0, 1e10, WithParallelism(4)),
			expect: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expect, test.stream.AnyMatch(func(item any) bool {
				return item.(int64) >= 500
			}))
		})
	}
}

func TestConcurrentStream_AllMatch(t *testing.T) {
	tests := []struct {
		name   string
		stream Stream
		expect bool
	}{
		{
			name:   "empty stream with no parallelism",
			stream: Just([]any{}, WithSync()),
			expect: true,
		},
		{
			name:   "empty stream with parallelism",
			stream: Just([]any{}, WithParallelism(4)),
			expect: true,
		},
		{
			name:   "non-empty stream with no parallelism and not all match",
			stream: Range(400, 1e10, WithSync()),
			expect: false,
		},
		{
			name:   "non-empty stream with parallelism and not all match",
			stream: Range(400, 1e10, WithParallelism(4)),
			expect: false,
		},
		{
			name:   "non-empty stream with no parallelism and all match",
			stream: Range(500, 1000, WithSync()),
			expect: true,
		},
		{
			name:   "non-empty stream with parallelism and all match",
			stream: Range(500, 1000, WithParallelism(4)),
			expect: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expect, test.stream.AllMatch(func(item any) bool {
				return item.(int64) >= 500
			}))
		})
	}
}

func TestConcurrentStream_NoneMatch(t *testing.T) {
	tests := []struct {
		name   string
		stream Stream
		expect bool
	}{
		{
			name:   "empty stream with no parallelism",
			stream: Just([]any{}, WithSync()),
			expect: true,
		},
		{
			name:   "empty stream with parallelism",
			stream: Just([]any{}, WithParallelism(4)),
			expect: true,
		},
		{
			name:   "non-empty stream with no parallelism and partial match",
			stream: Range(400, 1e10, WithSync()),
			expect: false,
		},
		{
			name:   "non-empty stream with parallelism and and partial match",
			stream: Range(400, 1e10, WithParallelism(4)),
			expect: false,
		},
		{
			name:   "non-empty stream with no parallelism and all match",
			stream: Range(500, 1000, WithSync()),
			expect: false,
		},
		{
			name:   "non-empty stream with parallelism and all match",
			stream: Range(500, 1000, WithParallelism(4)),
			expect: false,
		},
		{
			name:   "non-empty stream with no parallelism and no match",
			stream: Range(0, 500, WithSync()),
			expect: true,
		},
		{
			name:   "non-empty stream with parallelism and no match",
			stream: Range(0, 500, WithParallelism(4)),
			expect: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expect, test.stream.NoneMatch(func(item any) bool {
				return item.(int64) >= 500
			}))
		})
	}
}

func TestConcurrentStream_FindFirst(t *testing.T) {
	tests := []struct {
		name        string
		stream      Stream
		expect      any
		expectFound bool
	}{
		{
			name:        "empty stream with no parallelism",
			stream:      Just([]any{}, WithSync()),
			expect:      nil,
			expectFound: false,
		},
		{
			name:        "empty stream with parallelism",
			stream:      Just([]any{}, WithParallelism(4)),
			expect:      nil,
			expectFound: false,
		},
		{
			name:        "non-empty stream with no parallelism",
			stream:      Range(0, 1000, WithSync()),
			expect:      int64(500),
			expectFound: true,
		},
		{
			name:        "non-empty stream with parallelism",
			stream:      Range(0, 1000, WithParallelism(4)),
			expect:      int64(500),
			expectFound: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualItem, found := test.stream.Filter(func(item any) bool {
				return item.(int64) >= 500
			}).FindFirst()
			require.Equal(t, test.expectFound, found)
			if test.expectFound {
				require.GreaterOrEqual(t, actualItem, test.expect)
			}
		})
	}
}

func TestConcurrentStream_Count(t *testing.T) {
	tests := []struct {
		name   string
		stream Stream
		expect int64
	}{
		{
			name:   "empty stream with no parallelism",
			stream: Just([]any{}, WithSync()),
			expect: 0,
		},
		{
			name:   "empty stream with parallelism",
			stream: Just([]any{}, WithParallelism(4)),
			expect: 0,
		},
		{
			name:   "non-empty stream with no parallelism",
			stream: Range(0, 1000, WithSync()),
			expect: 1000,
		},
		{
			name:   "non-empty stream with parallelism",
			stream: Range(0, 1000, WithParallelism(4)),
			expect: 1000,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expect, test.stream.Count())
		})
	}
}

func TestConcurrentStream_Reduce(t *testing.T) {
	tests := []struct {
		name   string
		stream Stream
		expect any
	}{
		{
			name:   "empty stream with no parallelism",
			stream: Just([]any{}, WithSync()),
			expect: int64(0),
		},
		{
			name:   "empty stream with parallelism",
			stream: Just([]any{}, WithParallelism(4)),
			expect: int64(0),
		},
		{
			name:   "non-empty stream with no parallelism",
			stream: Range(0, 1000, WithSync()),
			expect: int64(499500),
		},
		{
			name:   "non-empty stream with parallelism",
			stream: Range(0, 1000, WithParallelism(4)),
			expect: int64(499500),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expect, test.stream.Reduce(int64(0), func(a, b any) any {
				return a.(int64) + b.(int64)
			}, WithSync()))
		})
	}
}
