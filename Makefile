build:
	go build -o bin/ova-conversation-api cmd/ova-conversation-api/main.go

run:
	go run cmd/ova-conversation-api/main.go

deps:
	@[ -f go.mod ] || go mod init github.com//ozonva/ova-conversation-api
	find . -name go.mod -exec bash -c 'pushd "$${1%go.mod}" && go mod tidy && popd' _ {} \;

bin-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go install github.com/envoyproxy/protoc-gen-validate@$(PGV_VERSION)

generate:
	protoc \
	--go_out=pkg/api --go_opt=paths=import \
	--go-grpc_out=pkg/api --go-grpc_opt=paths=import \
	api/api.proto