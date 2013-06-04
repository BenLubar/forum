package main

import (
	"testing"
)

func TestDiscussionKey(t *testing.T) {
	for _, test := range []struct {
		in  uint64
		out string
	}{
		{0, "discussion_0000000000000000"},
		{1, "discussion_0000000000000001"},
		{0x10000, "discussion_0000000000010000"},
		{0x123456789abcdef, "discussion_0123456789abcdef"},
		{^uint64(0), "discussion_ffffffffffffffff"},
		{1 << 63, "discussion_8000000000000000"},
	} {
		got := DiscussionKey(test.in)
		if want := test.out; got != want {
			t.Errorf("for %016x want %q got %q", test.in, test.out, got)
		}
	}
}

func BenchmarkDiscussionKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DiscussionKey(100)
	}
}
