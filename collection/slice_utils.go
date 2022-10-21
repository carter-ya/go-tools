package collection

import "math/rand"

// Reverse reverses the slice in place.
func Reverse[T any](s []T) {
	half := len(s) / 2
	for i := 0; i < half; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
}

// Shuffle shuffles the slice in place.
// See https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle.
func Shuffle[T any](s []T) {
	for i := len(s) - 1; i > 0; i-- {
		idx := rand.Intn(i + 1)
		s[i], s[idx] = s[idx], s[i]
	}
}
