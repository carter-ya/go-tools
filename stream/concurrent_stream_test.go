package stream

import (
	"github.com/carter-ya/go-tools/collection"
	"github.com/stretchr/testify/require"
	"sort"
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
