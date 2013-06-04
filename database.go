package main

import (
	"log"
	"sync"

	"github.com/couchbaselabs/go-couchbase"
	"github.com/dustin/gomemcached"
)

const (
	db_ddoc_version_key = "/@ddocVersion"
	db_ddoc_name        = "forum"
)

var db_init_once sync.Once

func db_init() {
	var ddocVersion uint64
	err := Bucket.Get(db_ddoc_version_key, &ddocVersion)
	if err != nil && !gomemcached.IsNotFound(err) {
		log.Fatalf("fatal error getting ddoc version from database: %v", err)
	}
	if wanted := uint64(1); ddocVersion != wanted {
		log.Printf("updating ddoc from version %d to %d", ddocVersion, wanted)
		err = Bucket.PutDDoc(db_ddoc_name, couchbase.DDocJSON{
			Views: map[string]couchbase.ViewDefinition{
				"discussion-posts": {
					Map: `function( doc, meta ) {
	if ( doc.Type == "post" ) {
		emit( doc.Discussion, null );
	}
}`,
					Reduce: `_count`,
				},

				"updated-discussions": {
					Map: `function( doc, meta ) {
	if ( doc.Type == "discussion" ) {
		emit( dateToArray( doc.Modified ), null );
	}
}`,
				},
			},
		})
		if err != nil {
			log.Fatalf("fatal error storing ddoc: %v", err)
		}
		err = Bucket.Set(db_ddoc_version_key, 0, wanted)
		if err != nil {
			log.Fatalf("fatal error storing ddoc version: %v", err)
		}
	}
}
