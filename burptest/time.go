package burptest

import (
	"math/rand"
	"time"
)

func RandTime() time.Time {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000
	return time.Unix(randomTime, 0).UTC()
}
