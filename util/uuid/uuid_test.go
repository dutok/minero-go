package uuid

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	i := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	e := "01020304-0506-0708-090a-0b0c0d0e0f10"
	o := strings.ToLower(UUID(i).String())

	if o != e {
		t.Fatalf("Expecting %q, got %q.", e, o)
	}
}

func BenchmarkUUID4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = UUID4().String()
	}
}
