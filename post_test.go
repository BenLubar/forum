package main

import (
	"html/template"
	"testing"
	"time"
)

func TestPost(t *testing.T) {
	p1, err1 := NewPost(2, 42)
	p2, err2 := NewPost(3, 42)

	if err1 != nil {
		t.Fatalf("NewPost failed: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("NewPost failed: %v", err2)
	}

	if want, got := post_type, p1.Type; want != got {
		t.Errorf("for Type want %q got %q", want, got)
	}
	if want, got := post_type, p2.Type; want != got {
		t.Errorf("for Type want %q got %q", want, got)
	}

	if want, got := uint64(2), p1.Author; want != got {
		t.Errorf("for Author want %d got %d", want, got)
	}
	if want, got := uint64(3), p2.Author; want != got {
		t.Errorf("for Author want %d got %d", want, got)
	}

	if want, got := template.HTML(""), p1.Content; want != got {
		t.Errorf("for Content want %q got %q", want, got)
	}
	if want, got := template.HTML(""), p2.Content; want != got {
		t.Errorf("for Content want %q got %q", want, got)
	}

	if want, got := uint64(42), p1.Discussion; want != got {
		t.Errorf("for Discussion want %d got %d", want, got)
	}
	if want, got := uint64(42), p2.Discussion; want != got {
		t.Errorf("for Discussion want %d got %d", want, got)
	}

	if want, got := uint64(1), p1.ID; want != got {
		t.Errorf("for ID want %d got %d", want, got)
	}
	if want, got := uint64(2), p2.ID; want != got {
		t.Errorf("for ID want %d got %d", want, got)
	}

	if want, got := 5*time.Second, time.Since(p1.Modified); want < got || got < 0 {
		t.Errorf("for Modified want %v got %v", want, got)
	}
	if want, got := 5*time.Second, time.Since(p2.Modified); want < got || got < 0 {
		t.Errorf("for Modified want %v got %v", want, got)
	}

	if !p1.Modified.Before(p2.Modified) {
		t.Error("modified 1 ≤ modified 2")
	}

	p1get, err := GetPost(p1.ID)
	if err == nil {
		t.Errorf("GetPost succeeded when it shouldn't have: %+v", p1get)
	}

	err = UpdateDiscussion(&Discussion{
		Type: discussion_type,
		ID:   42,
	})
	if err != nil {
		t.Fatalf("UpdateDiscussion failed: %v", err)
	}

	err = UpdatePost(p2)
	if err != nil {
		t.Fatalf("UpdatePost failed: %v", err)
	}

	p1get, err = GetPost(p1.ID)
	if err == nil {
		t.Errorf("GetPost succeeded when it shouldn't have: %+v", p1get)
	}

	p2get, err := GetPost(p2.ID)
	if err != nil {
		t.Errorf("GetPost failed when it shouldn't have: %v", err)
	}

	if *p2 != *p2get {
		t.Errorf("Posts not equal: %+v != %+v", *p2, *p2get)
	}
}

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
