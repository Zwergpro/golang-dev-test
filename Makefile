.PHONY: run deps migration

LOCAL_BIN:=$(CURDIR)/bin
GO:=$(GOROOT)/bin/go


run:
	$(GO) run cmd/bot/main.go


deps:
	GOBIN=$(LOCAL_BIN) $(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
	GOBIN=$(LOCAL_BIN) $(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
	GOBIN=$(LOCAL_BIN) $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go && \
	GOBIN=$(LOCAL_BIN) $(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc


MIGRATIONS_DIR=./migrations
migration:
	goose -dir=${MIGRATIONS_DIR} create $(NAME) sql

migrate:
	goose -dir=${MIGRATIONS_DIR} postgres "postgres://postgres:postgres@localhost:6432/postgres?sslmode=disable" up
