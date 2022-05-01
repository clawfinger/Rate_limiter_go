generate:
	rm -rf api/generated
	mkdir -p api/generated

	protoc \
		--proto_path=api/ \
		--go_out=api/generated \
		--go-grpc_out=api/generated \
		api/*.proto

build:
	go build -v -o ./bin ./cmd/limiter

db:
	docker start redis || docker run -p 6379:6379 -d redis

test:
	go test ./tests -race -count 10

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	$(shell go env GOPATH)/bin/golangci-lint run ./...

.PHONY: generate build