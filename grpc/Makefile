.PHONY: version generate run

version:
	go version
	protoc --version
	grpcurl --version

generate:
	protoc -I . \
	--go_out ./pkg/api/product --go_opt paths=source_relative \
	--go-grpc_out ./pkg/api/product --go-grpc_opt paths=source_relative \
	./product.proto

run-all:
	go run ./cmd/run