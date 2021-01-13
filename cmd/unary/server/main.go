package main

import (
	"flag"
	"github.com/ppzxc/go-grpc-examples-benchmark/cmd/unary"
	"runtime"
)

var port = flag.Int("port", 9990, "The server port")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	unary.ExecGRPCServer(*port)
}
