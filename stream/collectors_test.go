package stream

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCollector_Slice(t *testing.T) {
	expect := []int64{1, 2, 3, 4}
	require.Equal(t, expect, Just[int64](expect).Collect(SliceSupplier[int64](), SliceAccumulator[int64]()))
}

func TestCollector_Joining(t *testing.T) {
	expect := []int64{1, 2, 3, 4}
	require.Equal(t, "1,2,3,4", Just[int64](expect).Collect(JoiningSupplier[int64](), JoiningAccumulator[int64](",")))

	expect = []int64{}
	require.Equal(t, "", Just[int64](expect).Collect(JoiningSupplier[int64](), JoiningAccumulator[int64](",")))

	expect2 := []s{
		{S: "a", I: 1},
		{S: "b", I: 2},
	}
	require.Equal(t, "a:1,b:2", Just[s](expect2).Collect(JoiningSupplier[s](), JoiningAccumulator[s](",")))
}

var _ fmt.Stringer = (*s)(nil)

type s struct {
	S string
	I int
}

func (s s) String() string {
	return fmt.Sprintf("%s:%d", s.S, s.I)
}
