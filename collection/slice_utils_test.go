package collection

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChunk(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		size   int
		expect [][]int
	}{
		{
			name:   "empty",
			input:  []int{},
			size:   2,
			expect: [][]int{},
		},
		{
			name:  "size 1",
			input: []int{1, 2, 3, 4},
			size:  1,
			expect: [][]int{
				{1},
				{2},
				{3},
				{4},
			},
		},
		{
			name:  "size 2",
			input: []int{1, 2, 3, 4},
			size:  2,
			expect: [][]int{
				{1, 2},
				{3, 4},
			},
		},
		{
			name:  "size 3",
			input: []int{1, 2, 3, 4},
			size:  3,
			expect: [][]int{
				{1, 2, 3},
				{4},
			},
		},
		{
			name:  "size 4",
			input: []int{1, 2, 3, 4},
			size:  4,
			expect: [][]int{
				{1, 2, 3, 4},
			},
		},
		{
			name:  "size 5",
			input: []int{1, 2, 3, 4},
			size:  5,
			expect: [][]int{
				{1, 2, 3, 4},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expect, Chunk[int](test.input, test.size))
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		expect []int
	}{
		{
			name:   "empty",
			input:  []int{},
			expect: []int{},
		},
		{
			name:   "single",
			input:  []int{1},
			expect: []int{1},
		},
		{
			name:   "more",
			input:  []int{1, 2, 3, 4},
			expect: []int{4, 3, 2, 1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Reverse[int](test.input)
			require.Equal(t, test.expect, test.input)
		})
	}
}
