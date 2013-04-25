package auth

import (
	"testing"
)

// Test values taken from:
// http://wiki.vg/Protocol_Encryption
var tests = []struct {
	in, out string
}{
	{"Notch", "4ed1f46bbe04bc756bcb17c0c7ce3e4632f06a48"},
	{"jeb_", "-7c9d5b0044c130109a5d7b5fb5c317c02b4e28c1"},
	{"simon", "88e16a1019277b15d58faf0541e11910eb756f6"},
}

func TestAuthDigest(t *testing.T) {
	for index, tt := range tests {
		got := AuthDigest([]byte(tt.in))
		if got != tt.out {
			t.Fatalf("%d. %s failed!\nExpected: %x\nGot: %x", index, tt.out, got)
		}
	}
}
