package main

import (
	"context"
	"flag"
	pb "github.com/ppzxc/go-grpc-examples-benchmark/proto/example"
	"github.com/ppzxc/go-grpc-examples-benchmark/utils"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
)

var (
	port          = flag.Int("port", 9990, "The server port")
	payloadLength = flag.Int("len", 4*1023*1024, "The send payload length")
	payload       = []byte("")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	payload = utils.GenerateRandomBytes(*payloadLength)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Print(err)
	}

	s := grpc.NewServer()
	pb.RegisterExampleServer(s, &server{})

	log.Printf("gRPC BiStreamServer [%s]", l.Addr().String())

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

func (s *server) BiStream(stream pb.Example_BiStreamServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)

	var cn int32
	var wn int32

	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			req, err := stream.Recv()
			if err != nil {
				log.Printf("BiStreamServer, stream.Recv() error occurred %v\r\n", err)
				if err == io.EOF {
					return
				}
				return
			}
			if cn == 0 && wn == 0 {
				cn = req.GetConnNumber()
				wn = req.GetWorkerNumber()
			}
			log.Printf("CN:%d WN:%d RECEIVE RESPONSE UID:%d MLEN:%d LEN:%d\r\n", cn, wn, req.GetUid(), len(req.GetMessage()), req.GetLen())
		}
	}()

	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			serverSend := &pb.Response{
				Uid:     utils.Uuid.Get(),
				Message: payload,
				Len:     int32(len(payload)),
			}

			err := stream.Send(serverSend)
			if err != nil {
				log.Printf("BiStreamServer, stream.Send() error occurred %v\r\n", err)
				return
			}
			log.Printf("CN:%d WN:%d SUBMIT REQUEST UID:%d MLEN:%d LEN:%d\r\n", cn, wn, serverSend.GetUid(), len(serverSend.GetMessage()), serverSend.GetLen())
		}
	}()
	wg.Wait()
	return nil
}
