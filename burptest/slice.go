package burptest

import "math/rand"

// SliceItem picks up a random item from given slice.
func SliceItem[T any](slice []T) T {
	return slice[rand.Intn(len(slice)-1)]
}
