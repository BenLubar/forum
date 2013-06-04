package main

import (
	"testing"
)

func TestPostKey(t *testing.T) {
	for _, test := range []struct {
		in  uint64
		out string
	}{
		{0, "post_0000000000000000"},
		{1, "post_0000000000000001"},
		{0x10000, "post_0000000000010000"},
		{0x123456789abcdef, "post_0123456789abcdef"},
		{^uint64(0), "post_ffffffffffffffff"},
		{1 << 63, "post_8000000000000000"},
	} {
		got := PostKey(test.in)
		if want := test.out; got != want {
			t.Errorf("for %016x want %q got %q", test.in, test.out, got)
		}
	}
}

func BenchmarkPostKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = PostKey(100)
	}
}
