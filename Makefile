.PHONY: run deps

LOCAL_BIN:=$(CURDIR)/bin
GO:=$(GOROOT)/bin/go


run:
	$(GO) run cmd/bot/main.go


deps:
	GOBIN=$(LOCAL_BIN) $(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
	GOBIN=$(LOCAL_BIN) $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go && \
	GOBIN=$(LOCAL_BIN) $(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc