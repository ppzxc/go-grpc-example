package main

import (
	"context"
	pb "github.com/ppzxc/go-grpc-examples-benchmark/proto/example"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

type server struct {
	pb.UnimplementedExampleServer
}

func (u *server) Echo(ctx context.Context, Message *pb.Request) (*pb.Response, error) {
	return &pb.Response{
		Uid:     Message.Uid,
		Message: Message.Message,
		Len:     int32(len(Message.Message)),
	}, nil
}

func init() {
	go func() {
		l, err := net.Listen("tcp", ":9999")
		if err != nil {
			log.Print(err)
		}

		s := grpc.NewServer()
		pb.RegisterExampleServer(s, &server{})
		if err := s.Serve(l); err != nil {
			log.Println(err)
		}
	}()
}

func BenchmarkLocalServer(b *testing.B) {
	conn, err := grpc.Dial("localhost:9999", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	c := pb.NewExampleClient(conn)
	payload := GenerateRandomBytes(*payloadLength)

	b.Run("SINGLE", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pur := &pb.Request{
				Uid:          uint64(i),
				Message:      payload,
				Len:          int32(len(payload)),
				ConnNumber:   int32(1),
				WorkerNumber: int32(1),
			}
			ur, err := c.Echo(context.Background(), pur)
			if err != nil {
				log.Fatal(err)
			}
			b.SetBytes(int64(len(payload)))
			_ = ur
		}
	})
}
