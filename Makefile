unary_client:
	go build -o ./unary_client cmd/example/unary/client/main.go

unary_server:
	go build -o ./unary_server cmd/example/unary/server/main.go

all:
	go build -o ./unary_client cmd/example/unary/client/main.go
	go build -o ./unary_server cmd/example/unary/server/main.go
