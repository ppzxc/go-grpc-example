package main

import (
	"context"
	"crypto/rand"
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

func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

var (
	ip            = flag.String("ip", "localhost", "The server port")
	port          = flag.Int("port", 9990, "The server port")
	payloadLength = flag.Int("len", 4*1023*1024, "The send payload length")
	connCount     = flag.Int("conn", 10, "The connection count to grpc server")
	workerCount   = flag.Int("worker", 10, "The conn per worker")
)

var connArray []*grpc.ClientConn

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	payload := GenerateRandomBytes(*payloadLength)

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
			NewGRPCClient(ctx, cc, i, conn, &wg, payload)
		}
		connArray = append(connArray, conn)
	}
	wg.Wait()
}

func NewGRPCClient(ctx context.Context, cn int, gn int, conn *grpc.ClientConn, wg *sync.WaitGroup, payload []byte) {
	subCtx, cancel := context.WithCancel(ctx)
	c := pb.NewExampleClient(conn)
	var subWg sync.WaitGroup
	go func() {
		select {
		case <-subCtx.Done():
			break
		default:
		}
		bs, err := c.BiStream(subCtx)
		if err != nil {
			if err != nil {
				log.Fatal(err)
			}
		}

		subWg.Add(1)
		go func() {
			defer func() {
				cancel()
				subWg.Done()
			}()

			for {
				select {
				case <-subCtx.Done():
					return
				default:
				}
				pur := &pb.Request{
					Uid:          utils.Uuid.Get(),
					Message:      payload,
					Len:          int32(len(payload)),
					ConnNumber:   int32(cn),
					WorkerNumber: int32(gn),
				}
				err = bs.Send(pur)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("CN:%d WN:%d SUBMIT REQUEST UID:%d MLEN:%d LEN:%d\r\n", cn, gn, pur.GetUid(), len(pur.GetMessage()), pur.GetLen())
			}
		}()

		subWg.Add(1)
		go func() {
			defer func() {
				cancel()
				subWg.Done()
			}()

			for {
				select {
				case <-subCtx.Done():
					return
				default:
				}

				res, err := bs.Recv()
				if err != nil {
					log.Printf("BiStreamClient, stream.Recv() error occurred %v\r\n", err)
					if err == io.EOF {
						return
					}
					return
				}
				log.Printf("CN:%d WN:%d RECEIVE RESPONSE UID:%d MLEN:%d LEN:%d\r\n", cn, gn, res.GetUid(), len(res.GetMessage()), res.GetLen())
			}
		}()

		subWg.Wait()
		wg.Done()
	}()
}
