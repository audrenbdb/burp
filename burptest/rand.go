// Package burptest provides utility functions designed to ease tests.
package burptest

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
