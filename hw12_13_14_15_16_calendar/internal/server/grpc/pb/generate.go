package pb

//go:generate protoc --proto_path=../../../../api --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative calendar.proto

// This file enables go generate support.
// Run 'go generate ./internal/server/grpc/pb/...' to generate gRPC code.
