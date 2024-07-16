#/bin/bash
protoc --go_out=internal/remote/generated --go-grpc_out=internal/remote/generated proto/hello.proto
