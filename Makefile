echo_client:
	go build -o ./echo_client cmd/example/echo/client/main.go

echo_server:
	go build -o ./echo_server cmd/example/echo/server/main.go

clientStream_client:
	go build -o ./clientStream_client cmd/example/clientStream/client/main.go

clientStream_server:
	go build -o ./clientStream_server cmd/example/clientStream/server/main.go

serverStream_client:
	go build -o ./serverStream_client cmd/example/serverStream/client/main.go

serverStream_server:
	go build -o ./serverStream_server cmd/example/serverStream/server/main.go

biStream_client:
	go build -o ./biStream_client cmd/example/biStream/client/main.go

biStream_server:
	go build -o ./biStream_server cmd/example/biStream/server/main.go

all:
	go build -o ./echo_client cmd/example/echo/client/main.go
	go build -o ./echo_server cmd/example/echo/server/main.go
	go build -o ./clientStream_client cmd/example/clientStream/client/main.go
	go build -o ./clientStream_server cmd/example/clientStream/server/main.go
	go build -o ./serverStream_client cmd/example/serverStream/client/main.go
	go build -o ./serverStream_server cmd/example/serverStream/server/main.go
	go build -o ./biStream_client cmd/example/biStream/client/main.go
	go build -o ./biStream_server cmd/example/biStream/server/main.go
