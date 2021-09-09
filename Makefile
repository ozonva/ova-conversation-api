ifneq (,$(wildcard ./.env))
    include .env
    export
endif

DBSTRING:="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable"

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
	protoc -I vendor.protogen \
	--go_out=pkg/api --go_opt=paths=import \
	--go-grpc_out=pkg/api --go-grpc_opt=paths=import \
	--validate_out lang=go:pkg/api \
	api/api.proto

generate-vendor-proto:
	mkdir -p vendor.protogen
	mkdir -p vendor.protogen/api
	cp api/api.proto vendor.protogen/api/api.proto
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
		mkdir -p vendor.protogen/google/ &&\
		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
		rm -rf vendor.protogen/googleapis ;\
	fi
	@if [ ! -d vendor.protogen/github.com/envoyproxy ]; then \
		mkdir -p vendor.protogen/github.com/envoyproxy &&\
		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/github.com/envoyproxy/protoc-gen-validate ;\
	fi

migrate:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DBSTRING) goose -dir migration up
