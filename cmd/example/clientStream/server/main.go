package main

import (
	"flag"
	pb "github.com/ppzxc/go-grpc-examples-benchmark/proto/example"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

var port = flag.Int("port", 9990, "The server port")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	l, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Print(err)
	}

	s := grpc.NewServer()
	pb.RegisterExampleServer(s, &server{})

	log.Printf("gRPC ClientStreamServer [%s]", l.Addr().String())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		sig := <-sigs
		log.Println(sig.String())
		if err := l.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err := s.Serve(l); err != nil {
		log.Println(err)
	}
}

type server struct {
	pb.UnimplementedExampleServer
}

func (s *server) ClientStream(stream pb.Example_ClientStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			log.Printf("ClientStreamServer, stream.Recv() error occurred %v\r\n", err)
			if err == io.EOF {
				return nil
			}
			return err
		}
		log.Printf("ClientStreamServer, recv message send UID:%d CN:%d WN:%d LEN:%d \r\n", req.GetUid(), req.GetConnNumber(), req.GetWorkerNumber(), req.Len)
		if req.GetLen() == -1 {
			break
		}
	}

	if err := stream.SendAndClose(&pb.Response{
		Message: []byte("DONE"),
	}); err != nil {
		log.Printf("ClientStreamServer, stream.SendAndClose() error occurred %v\r\n", err)
	}
	return nil
}
