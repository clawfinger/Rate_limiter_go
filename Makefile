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

.PHONY: generate build