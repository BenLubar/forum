package main

import (
	"encoding/json"
	"time"
)

type Discussion struct {
	Type     string
	ID       uint64
	Title    string
	Modified time.Time
}

const (
	discussion_type       = "discussion"
	discussion_key_prefix = "discussion_"
	discussion_incr_key   = "incr_discussion"
)

// NewDiscussion creates a new discussion with a unique ID. Use UpdateDiscussion
// to save the discussion to the database.
func NewDiscussion() (*Discussion, error) {
	id, err := Bucket.Incr(discussion_incr_key, 1, 1, 0)
	if err != nil {
		return nil, err
	}

	return &Discussion{
		Type:     discussion_type,
		ID:       id,
		Modified: time.Now().UTC(),
	}, nil
}

// GetDiscussion retrieves a discussion from the database.
func GetDiscussion(id uint64) (*Discussion, error) {
	var d Discussion

	err := Bucket.Get(DiscussionKey(id), &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// UpdateDiscussion saves the given Discussion to the database.
func UpdateDiscussion(d *Discussion) error {
	d.Modified = time.Now().UTC()

	return Bucket.Set(DiscussionKey(d.ID), 0, d)
}

// TouchDiscussion updates the Modified timestamp on a discussion to the
// current time.
func TouchDiscussion(id uint64) error {
	return Bucket.Update(DiscussionKey(id), 0, func(in []byte) ([]byte, error) {
		var d Discussion
		err := json.Unmarshal(in, &d)
		if err != nil {
			return in, err
		}

		d.Modified = time.Now().UTC()

		return json.Marshal(&d)
	})
}

// DiscussionKey generates a database key for a given discussion ID, in the
// format "discussion_012345679abcdef".
func DiscussionKey(id uint64) string {
	const hexChars = "0123456789abcdef"

	key := make([]byte, len(discussion_key_prefix)+(64/4))
	copy(key, discussion_key_prefix)

	j := len(discussion_key_prefix)
	for i := 64 - 4; i >= 0; i -= 4 {
		key[j] = hexChars[(id>>uint(i))&0xf]
		j++
	}

	return string(key)
}
