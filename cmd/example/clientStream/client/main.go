package main

import (
	"context"
	"flag"
	pb "github.com/ppzxc/go-grpc-examples-benchmark/proto/example"
	"github.com/ppzxc/go-grpc-examples-benchmark/utils"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
)

var (
	ip            = flag.String("ip", "localhost", "The server port")
	port          = flag.Int("port", 9990, "The server port")
	payloadLength = flag.Int("len", 4*1023*1024, "The send payload length")
	connCount     = flag.Int("conn", 10, "The connection count to grpc server")
	workerCount   = flag.Int("worker", 10, "The conn per worker")
	sendingCount  = flag.Int("count", 100000, "The conn per worker")
)

var connArray []*grpc.ClientConn

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	payload := utils.GenerateRandomBytes(*payloadLength)

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
	c := pb.NewExampleClient(conn)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			break
		default:
		}
		cs, err := c.ClientStream(ctx)
		if err != nil {
			if err != nil {
				log.Fatal(err)
			}
		}
		for i := 0; i < *sendingCount; i++ {
			pur := &pb.Request{
				Uid:          utils.Uuid.Get(),
				Message:      payload,
				Len:          int32(len(payload)),
				ConnNumber:   int32(cn),
				WorkerNumber: int32(gn),
			}
			err = cs.Send(pur)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("CN:%d WN:%d SUBMIT REQUEST UID:%d MLEN:%d LEN:%d", cn, gn, pur.GetUid(), len(pur.GetMessage()), pur.GetLen())
		}
		err = cs.Send(&pb.Request{Len: int32(-1)})
		if err != nil {
			log.Fatal(err)
		}
		res, e := cs.CloseAndRecv()
		if e != nil {
			log.Fatal(e)
			return
		}
		log.Printf("%s\r\n", res.String())
	}()
}
