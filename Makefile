client:
	go build -o ./client cmd/unary/client/main.go

server:
	go build -o ./server cmd/unary/server/main.go

all:
	go build -o ./client cmd/unary/client/main.go
	go build -o ./server cmd/unary/server/main.go
