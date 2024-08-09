SOURCES=$(shell find . -name "*.go")

.PHONY: format
format:
	golangci-lint run --fix

.PHONY: check-format
check-format:
	golangci-lint run

.PHONY: build
build: pwstore_backend

.PHONY: codegen
codegen: internal/api/api.go

.PHONY: run
run: build
	./pwstore_backend

.PHONY: test
test: build
	go test -v ./test

pwstore_backend: ${SOURCES}
	go build -o pwstore_backend cmd/rest/main.go

internal/api/api.go: api/api.yaml configs/codegen/server.yaml
	mkdir -p internal/api
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config configs/codegen/server.yaml  api/api.yaml
