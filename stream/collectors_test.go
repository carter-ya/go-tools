package stream

import (
	"fmt"
	"github.com/carter-ya/go-tools/stream/collector"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCollector_Slice(t *testing.T) {
	expect := []int64{1, 2, 3, 4}
	require.Equal(t, expect, Just[int64](expect).Collect(
		NewToSliceCollector[int64](0),
	))
}

func TestCollector_Joining(t *testing.T) {
	expect := []int64{1, 2, 3, 4}
	require.Equal(t, "1,2,3,4", Just[int64](expect).Collect(
		NewJoiningCollector[int64](","),
	))

	expect = []int64{}
	require.Equal(t, "", Just[int64](expect).Collect(
		NewJoiningCollector[int64](","),
	))

	expect2 := []s{
		{S: "a", I: 1},
		{S: "b", I: 2},
	}
	require.Equal(t, "a:1,b:2", Just[s](expect2).Collect(
		NewJoiningCollector[s](",")),
	)
}

func TestCollector_GroupBy(t *testing.T) {
	expect := []int64{1, 2, 3, 4}
	require.Equal(t, map[bool][]int64{true: {1, 3}, false: {2, 4}}, Just[int64](expect).Collect(
		NewGroupByCollector(func(i int64) bool {
			return i%2 == 1
		}),
	))
}

func TestCollector_Count(t *testing.T) {
	expect := []int64{1, 2, 3, 4}
	require.Equal(t, len(expect), Just[int64](expect).Collect(
		NewCountCollector[int64](),
	))
}

func TestCollector_Sum(t *testing.T) {
	expect := []int64{1, 2, 3, 4}
	require.Equal(t, int64(10), Just[int64](expect).Collect(
		NewSumCollector[int64](collector.Identify[int64]()),
	))
}

func TestCollector_Avg(t *testing.T) {
	expect := []int64{1, 2, 3, 4}
	require.Equal(t, 2.5, Just[int64](expect).Collect(
		NewAvgCollector[int64, float64](func(item int64) float64 {
			return float64(item)
		}),
	))
}

func TestCollector_Max(t *testing.T) {
	expect := []int64{1, 2, 3, 4}
	require.Equal(t, int64(4), Just[int64](expect).Collect(
		NewMaxCollector[int64](collector.Identify[int64]()),
	))
}

func TestCollector_Min(t *testing.T) {
	expect := []int64{1, 2, 3, 4}
	require.Equal(t, int64(1), Just[int64](expect).Collect(
		NewMinCollector[int64](collector.Identify[int64]()),
	))
}

var _ fmt.Stringer = (*s)(nil)

type s struct {
	S string
	I int
}

func (s s) String() string {
	return fmt.Sprintf("%s:%d", s.S, s.I)
}
