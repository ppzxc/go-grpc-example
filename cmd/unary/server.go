package unary

import (
	"context"
	pb "github.com/ppzxc/go-grpc-examples-benchmark/proto/unary"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type UnaryEchoServer struct {
	Unary *pb.UnaryRequest
}

func (u *UnaryEchoServer) UnaryEcho(ctx context.Context, unaryMessage *pb.UnaryRequest) (*pb.UnaryResponse, error) {
	log.Printf("SUBMIT REQUEST CN:%d GN:%d UID:%d MLEN:%d LEN:%d", unaryMessage.GetConnNumber(), unaryMessage.GetWorkerNumber(), unaryMessage.GetUid(), len(unaryMessage.GetMessage()), unaryMessage.GetLen())
	return &pb.UnaryResponse{
		Uid:     unaryMessage.Uid,
		Message: unaryMessage.Message,
		Len:     int32(len(unaryMessage.Message)),
	}, nil
}

func ExecGRPCServer(p int) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(p))
	if err != nil {
		log.Print(err)
	}

	s := grpc.NewServer()
	pb.RegisterUnaryExampleServer(s, &UnaryEchoServer{})

	log.Printf("gRPC Server [%s]", l.Addr().String())

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
