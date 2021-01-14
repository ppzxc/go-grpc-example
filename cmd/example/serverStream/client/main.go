package main

import (
	"context"
	"flag"
	pb "github.com/ppzxc/go-grpc-examples-benchmark/proto/example"
	"github.com/ppzxc/go-grpc-examples-benchmark/utils"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
)

var (
	ip          = flag.String("ip", "localhost", "The server port")
	port        = flag.Int("port", 9990, "The server port")
	connCount   = flag.Int("conn", 10, "The connection count to grpc server")
	workerCount = flag.Int("worker", 10, "The conn per worker")
)

var connArray []*grpc.ClientConn

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-sigs
		log.Println(sig.String())
		for _, conn := range connArray {
			if err := conn.Close(); err != nil {
				log.Println(err)
			}
		}
		cancel()
	}()

	var wg sync.WaitGroup
	for cc := 0; cc < *connCount; cc++ {
		conn, err := grpc.Dial(*ip+":"+strconv.Itoa(*port), grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < *workerCount; i++ {
			wg.Add(1)
			NewGRPCClient(ctx, cc, i, conn, &wg)
		}
		connArray = append(connArray, conn)
	}
	wg.Wait()
}

func NewGRPCClient(ctx context.Context, cn int, gn int, conn *grpc.ClientConn, wg *sync.WaitGroup) {
	c := pb.NewExampleClient(conn)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				break
			default:
			}
			pur := &pb.Request{
				Uid:          utils.Uuid.Get(),
				Message:      []byte("HI"),
				Len:          int32(len([]byte("HI"))),
				ConnNumber:   int32(cn),
				WorkerNumber: int32(gn),
			}
			stream, err := c.ServerStream(ctx, pur)
			if err != nil {
				if err != nil {
					log.Fatal(err)
				}
			}
			for {
				in, err := stream.Recv()
				if err != nil {
					log.Printf("ServerStreamClient, stream.Recv() error occurred %v\r\n", err)
					if err == io.EOF {
						break
					}
					break
				}
				log.Printf("CN:%d WN:%d SUBMIT REQUEST UID:%d MLEN:%d LEN:%d ::==:: RECEIVED RESPONSE UID:%d MLEN:%d LEN:%d", cn, gn, pur.GetUid(), len(pur.GetMessage()), pur.GetLen(), in.GetUid(), len(in.GetMessage()), in.GetLen())
			}
		}
	}()
}
