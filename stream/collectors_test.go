package stream

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCollector_Slice(t *testing.T) {
	expect := []int64{1, 2, 3, 4}
	require.Equal(t, expect, Just[int64](expect).Collect(SliceSupplier[int64](), SliceAccumulator[int64]()))
}
