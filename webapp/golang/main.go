package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	bind string
	port int
)

func init() {
	flag.StringVar(&bind, "bind", "0.0.0.0", "bind address")
	flag.IntVar(&port, "port", 9292, "bind port")

	flag.Parse()
}

func main() {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", bind, port),
		Handler: serveMux(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
