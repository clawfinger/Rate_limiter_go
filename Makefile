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

.PHONY: generate build