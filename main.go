package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/couchbaselabs/go-couchbase"
)

var Bucket *couchbase.Bucket

func main() {
	var (
		addr = flag.String("addr", ":8015", "address to listen on")
		cert = flag.String("tls-cert", "", "cert.pem file")
		key  = flag.String("tls-key", "", "key.pem file")

		cburl    = flag.String("couchbase", "http://127.0.0.1:8091", "")
		cbpool   = flag.String("cb-pool", "default", "")
		cbbucket = flag.String("cb-bucket", "default", "")
	)

	flag.Parse()

	if (*cert == "") != (*key == "") {
		log.Fatalf("the -tls-cert flag requires the -tls-key flag, and vice-versa.")
	}

	b, err := couchbase.GetBucket(*cburl, *cbpool, *cbbucket)
	if err != nil {
		log.Fatalf("error connecting to Couchbase: %v", err)
	}
	Bucket = b
	go db_init_once.Do(db_init)

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("error listening on %q: %v", *addr, err)
	}

	if *cert != "" {
		config := &tls.Config{
			NextProtos:   []string{"http/1.1"},
			Certificates: make([]tls.Certificate, 1),
		}
		config.Certificates[0], err = tls.LoadX509KeyPair(*cert, *key)
		if err != nil {
			log.Fatalf("error loading TLS data: %v", err)
		}
		ln = tls.NewListener(ln, config)
	}

	log.Printf("Now listening on %v", ln.Addr())
	err = http.Serve(ln, nil)
	if err != nil {
		log.Fatalf("fatal error: %v", err)
	}
}
