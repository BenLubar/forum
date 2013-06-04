package main

import (
	"github.com/couchbaselabs/walrus"
)

func init() {
	Bucket = walrus.NewBucket("default")
	db_init_once.Do(db_init)
}
