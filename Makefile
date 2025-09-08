GO ?= go

GO_DEPS := \
	mvdan.cc/gofumpt@latest

.PHONY: all tidy run test fmt check-fmt vet generate tools deps install clean

all: build

test:
	$(GO) test ./...

fmt:
	gofumpt -l -w .

vet:
	$(GO) vet ./...

deps:
	$(GO) get $(GO_DEPS)
	$(GO) mod tidy

generate:
	$(GO) generate ./...
