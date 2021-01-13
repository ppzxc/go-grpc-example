package main

import (
	"context"
	"github.com/ppzxc/go-grpc-examples-benchmark/cmd/unary"
	pb "github.com/ppzxc/go-grpc-examples-benchmark/proto/unary"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func init() {
	go func() {
		l, err := net.Listen("tcp", ":9999")
		if err != nil {
			log.Print(err)
		}

		s := grpc.NewServer()
		pb.RegisterUnaryExampleServer(s, &unary.UnaryEchoServer{})
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
	c := pb.NewUnaryExampleClient(conn)
	payload := GenerateRandomBytes(*payloadLength)

	b.Run("SINGLE", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pur := &pb.UnaryRequest{
				Uid:          uint64(i),
				Message:      payload,
				Len:          int32(len(payload)),
				ConnNumber:   int32(1),
				WorkerNumber: int32(1),
			}
			ur, err := c.UnaryEcho(context.Background(), pur)
			if err != nil {
				log.Fatal(err)
			}
			b.SetBytes(int64(len(payload)))
			_ = ur
		}
	})
}
