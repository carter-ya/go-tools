package collection

import "math/rand"

// Shuffle shuffles the slice in place.
// See https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle.
func Shuffle[T any](s []T) {
	for i := len(s) - 1; i > 0; i-- {
		idx := rand.Intn(i + 1)
		s[i], s[idx] = s[idx], s[i]
	}
}
