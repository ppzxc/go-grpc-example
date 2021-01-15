package main

import (
	"flag"
	pb "github.com/ppzxc/go-grpc-examples-benchmark/proto/example"
	"github.com/ppzxc/go-grpc-examples-benchmark/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

var (
	port          = flag.Int("port", 9990, "The server port")
	count         = flag.Int("count", 10000, "The server port")
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

	log.Printf("gRPC ServerStreamServer [%s]", l.Addr().String())

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

func (s *server) ServerStream(req *pb.Request, stream pb.Example_ServerStreamServer) error {
	log.Printf("ServerStreamServer, recv message UID:%d MSG:%s CN:%d WN:%d LEN:%d SERVER PUSH START\r\n", req.GetUid(), string(req.GetMessage()), req.GetConnNumber(), req.GetWorkerNumber(), req.Len)
	for i := 0; i < *count; i++ {
		res := &pb.Response{
			Uid:          utils.Uuid.Get(),
			Message:      payload,
			Len:          int32(len(payload)),
			ConnNumber:   req.GetConnNumber(),
			WorkerNumber: req.GetWorkerNumber(),
		}
		if err := stream.Send(res); err != nil {
			log.Printf("ServerStreamServer, stream.Send() error occurred %v\r\n", err)
			return err
		}
		log.Printf("ServerStreamServer, push UID:%d CN:%d WN:%d LEN:%d\r\n", res.GetUid(), res.GetConnNumber(), res.GetWorkerNumber(), res.Len)
	}

	log.Printf("ServerStreamServer, Close Server Stream\r\n")
	if err := stream.Send(&pb.Response{
		Len: int32(-1),
	}); err != nil {
		log.Printf("ServerStreamServer, stream.Send() error occurred %v\r\n", err)
		return err
	}
	return nil
}
