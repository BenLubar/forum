package main

import (
	"html/template"
	"time"
)

type Post struct {
	Type       string
	ID         uint64
	Author     uint64
	Discussion uint64
	Content    template.HTML
	Modified   time.Time
}

const (
	post_type       = "post"
	post_key_prefix = "post_"
	post_incr_key   = "incr_post"
)

// NewPost creates a new post with a unique ID and the given author and
// discussion. Use UpdatePost to save the post to the database.
func NewPost(author, discussion uint64) (*Post, error) {
	db_init_once.Do(db_init)

	id, err := Bucket.Incr(post_incr_key, 1, 0, 0)
	if err != nil {
		return nil, err
	}

	return &Post{
		Type:       post_type,
		ID:         id,
		Author:     author,
		Discussion: discussion,
		Modified:   time.Now().UTC(),
	}, nil
}

// GetPost retrieves a single post from the database.
func GetPost(id uint64) (*Post, error) {
	db_init_once.Do(db_init)

	var p Post

	err := Bucket.Get(PostKey(id), &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// UpdatePost saves the given Post to the database and calls TouchDiscussion
// with the post's discussion ID.
func UpdatePost(p *Post) error {
	db_init_once.Do(db_init)

	p.Modified = time.Now().UTC()

	key := PostKey(p.ID)

	err := Bucket.Set(key, 0, p)
	if err != nil {
		return err
	}

	return TouchDiscussion(p.Discussion)
}

// PostKey generates a database key for a given post ID, in the format
// "post_012345679abcdef".
func PostKey(id uint64) string {
	const hexChars = "0123456789abcdef"

	key := make([]byte, len(post_key_prefix)+(64/4))
	copy(key, post_key_prefix)

	j := len(post_key_prefix)
	for i := 64 - 4; i >= 0; i -= 4 {
		key[j] = hexChars[(id>>uint(i))&0xf]
		j++
	}

	return string(key)
}
