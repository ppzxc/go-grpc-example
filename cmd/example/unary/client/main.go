package main

import (
	"context"
	"crypto/rand"
	"flag"
	pb "github.com/ppzxc/go-grpc-examples-benchmark/proto/example"
	"google.golang.org/grpc"
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
	payloadLength = flag.Int("len", 32768, "The send payload length")
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

type UUID struct {
	mu   sync.Mutex
	uuid uint64
}

func (u *UUID) Get() uint64 {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.uuid = u.uuid + 1
	return u.uuid
}

var Uuid UUID

func NewGRPCClient(ctx context.Context, cn int, gn int, conn *grpc.ClientConn, wg *sync.WaitGroup, payload []byte) {
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
				Uid:          Uuid.Get(),
				Message:      payload,
				Len:          int32(len(payload)),
				ConnNumber:   int32(cn),
				WorkerNumber: int32(gn),
			}
			ur, err := c.UnaryEcho(ctx, pur)

			if err != nil {
				log.Fatal(err)
			}
			log.Printf("CN:%d WN:%d SUBMIT REQUEST UID:%d MLEN:%d LEN:%d ::==:: RECEIVED RESPONSE UID:%d MLEN:%d LEN:%d", cn, gn, pur.GetUid(), len(pur.GetMessage()), pur.GetLen(), ur.GetUid(), len(ur.GetMessage()), ur.GetLen())

		}
	}()
}
