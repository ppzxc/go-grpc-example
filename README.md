# simple throughput check on real machine using grpc
- using grpc/grpc-go
- https://github.com/grpc/grpc-go

## usage
1. git clone https://github.com/ppzxc/go-grpc-examples-benchmark.git
2. cd go-grpc-examples-benchmark
3. make all
4. server side
 - ./echo_server -port 9990
5. client side
 - ./echo_client -ip 192.168.0.65 -port 9990 -len 65536 -conn 20 -worker 5
 - ./echo_client -ip 192.168.0.65 -port 9990 -len 32768 -conn 5 -worker 5
 - ./echo_client -ip 192.168.0.65 -port 9990 -len 512 -conn 5 -worker 1

## flags

```protobuf
message Request {
  uint64 uid = 1;
  bytes message = 2;
  int32 len = 3;
  int32 connNumber = 4;
  int32 workerNumber = 5;
}
```

 - ip
 - port
 - len => protobuf bytes message size
 - conn => connection count
 - worker => connection per worker count

## grpc/grpc-go
```
D:\go\go-grpc-examples-benchmark\proto\unary>protoc --version
libprotoc 3.14.0
D:\go\go-grpc-examples-benchmark\proto\unary>protoc --go_out=plugins=grpc:. *.proto
```

## reference
- https://github.com/grpc/grpc-go
- https://github.com/gogo/protobuf
- https://github.com/gogo/grpc-example

## todo
- stream

---------------
---------------
---------------
![du](./1.JPG)