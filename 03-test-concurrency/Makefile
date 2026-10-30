.PHONY: version generate run-all test test-race test-api

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

# Обычный прогон юнит-тестов
test:
	go test ./...

# Прогон всех тестов с race-детектором
test-race:
	go test -race ./...

# Интеграционные тесты
# -count=1 отключает кэширование, чтобы тесты реально ходили в сеть
test-api:
	go test -count=1 -tags=api ./test/...