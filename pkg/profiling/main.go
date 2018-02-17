// Package profiling
// Talk: https://www.youtube.com/watch?v=2h_NFBFrciI
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

// profile runtime  ------
// $> time go fmt std
// $> time -v go fmt std
// $> time -l go fmt std

// profile go build ------
// $> go build -toolexec="/usr/bin/time" .

// GODEBUG  ------
// $> env GODEBUG=gctrace=1 godoc -http:8080

// Profilers  ------
// before and after with pprof
// CPU and Memory, Block profiling - one kind of at a time

// enable debugging over http

func main() {
	log.Println(http.ListenAndServe("localhost:6060", nil))
}

// more here https://golang.org/pkg/net/http/pprof/ and https://blog.golang.org/profiling-go-programs
// $> go tool pprof http://localhost:6060/debug/pprof/heap
// $> go tool pprof http://localhost:6060/debug/pprof/profile
// graphical representations

// Framepointers supports perf
// go build -toolexec="perf stat" .
// flamegraph with go-torch

// go tool trace
// goroutine creation started and end
