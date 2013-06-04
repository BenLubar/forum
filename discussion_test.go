package main

import (
	"testing"
	"time"
)

func TestDiscussion(t *testing.T) {
	d, err := NewDiscussion()
	if err != nil {
		t.Fatalf("NewDiscussion() failed: %v", err)
	}

	if want, got := discussion_type, d.Type; want != got {
		t.Errorf("for Type want %q got %q", want, got)
	}
	if want, got := uint64(1), d.ID; want != got {
		t.Errorf("for ID want %d got %d", want, got)
	}
	if want, got := 5*time.Second, time.Since(d.Modified); want < got || got < 0 {
		t.Errorf("for Modified want %v got %v", want, got)
	}
	if want, got := "", d.Title; want != got {
		t.Errorf("for Title want %q got %q", want, got)
	}

	d2, err := GetDiscussion(d.ID)
	if err == nil {
		t.Errorf("GetDiscussion succeeded when it shouldn't have: %v", d2)
	}

	if err = TouchDiscussion(d.ID); err == nil {
		t.Error("TouchDiscussion succeeded when it shouldn't have")
	}

	if err = TouchDiscussion(d.ID); err == nil {
		t.Error("TouchDiscussion succeeded when it shouldn't have")
	}

	d2, err = GetDiscussion(d.ID)
	if err == nil {
		t.Errorf("GetDiscussion succeeded when it shouldn't have: %v", d2)
	}

	if err = UpdateDiscussion(d); err != nil {
		t.Fatalf("UpdateDiscussion() failed: %v", err)
	}

	d2, err = GetDiscussion(d.ID)
	if err != nil {
		t.Errorf("GetDiscussion failed when it shouldn't have: %v", err)
	}

	if *d != *d2 {
		t.Errorf("Discussions not equal: %+v != %+v", *d, *d2)
	}

	if err = TouchDiscussion(d.ID); err != nil {
		t.Errorf("TouchDiscussion failed when it shouldn't have: %v", err)
	}
}

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
