package sleep

import (
	"math/rand"
	"time"
)

func Random() {
	ms := rand.Intn(50) * 100
	time.Sleep(time.Millisecond * time.Duration(ms))
}
