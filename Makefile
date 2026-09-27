GO ?= go
PUBLIC_MODULE ?= github.com/example/your-service-go
PUBLIC_PACKAGES ?= .,httpclient

.DEFAULT_GOAL := check
.PHONY: check build test lint fmt api-compatibility

check: lint build test

build:
	$(GO) build ./...

test:
	$(GO) test -race ./...

lint:
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

api-compatibility:
	$(GO) run ./tools/compatibility -policy report -base previous-release -module "$(PUBLIC_MODULE)" -packages "$(PUBLIC_PACKAGES)"
