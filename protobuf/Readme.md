# How to generate the Go files from .proto

1. cd into proto-files
2. run this command

protoc --go_out=../generated --go_opt=paths=source_relative \
    --go-grpc_out=../generated --go-grpc_opt=paths=source_relative \
    ./dex.proto


