GO ?= go
GOLANGCI_LINT ?= golangci-lint
PUBLIC_MODULE ?= github.com/example/your-service-go
PUBLIC_PACKAGES ?= .,httpclient

.DEFAULT_GOAL := check
.PHONY: check build test lint fmt schema-examples api-compatibility

check: lint build test schema-examples

build:
	$(GO) build ./...

test:
	$(GO) test -race ./...

schema-examples:
	$(GO) run ./tools/schemaexamples

lint:
	$(GO) vet ./...
	$(GOLANGCI_LINT) run --timeout=5m ./...

fmt:
	$(GO) fmt ./...

api-compatibility:
	$(GO) run ./tools/compatibility -policy report -base previous-release -module "$(PUBLIC_MODULE)" -packages "$(PUBLIC_PACKAGES)"
