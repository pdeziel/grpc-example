package pb

//go:generate protoc -I../proto --go_opt=module=github.com/pdeziel/grpc-example/pb --go_out=. --go-grpc_opt=module=github.com/pdeziel/grpc-example/pb --go-grpc_out=. hello.proto
