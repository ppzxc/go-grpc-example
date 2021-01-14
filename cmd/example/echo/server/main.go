package main

import (
	"context"
	"flag"
	pb "github.com/ppzxc/go-grpc-examples-benchmark/proto/example"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

var port = flag.Int("port", 9990, "The server port")

type server struct {
	pb.UnimplementedExampleServer
}

func (u *server) Echo(ctx context.Context, Message *pb.Request) (*pb.Response, error) {
	log.Printf("SUBMIT REQUEST CN:%d GN:%d UID:%d MLEN:%d LEN:%d", Message.GetConnNumber(), Message.GetWorkerNumber(), Message.GetUid(), len(Message.GetMessage()), Message.GetLen())
	return &pb.Response{
		Uid:     Message.Uid,
		Message: Message.Message,
		Len:     int32(len(Message.Message)),
	}, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	l, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Print(err)
	}

	s := grpc.NewServer()
	pb.RegisterExampleServer(s, &server{})

	log.Printf("gRPC UnaryServer [%s]", l.Addr().String())

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
