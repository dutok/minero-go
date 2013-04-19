package simplex

import (
	"math"
	"testing"
	"testing/quick"
)

func BenchmarkFastFloor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fastfloor(0.5)
	}
}

func BenchmarkMathFloor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = int(math.Floor(0.5))
	}
}

func TestFloor(t *testing.T) {
	f := func(f float64) bool {
		return fastfloor(f) == int(math.Floor(f))
	}

	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}
