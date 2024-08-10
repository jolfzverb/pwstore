SOURCES=$(shell find internal cmd -name "*.go" -name "*.sql")

.PHONY: format
format:
	golangci-lint run --fix

.PHONY: check-format
check-format:
	golangci-lint run

.PHONY: build
build: pwstore_backend

.PHONY: codegen
codegen: internal/api/api.gen.go internal/clients/google_open_id/google_open_id.gen.go

.PHONY: run
run: build
	./pwstore_backend

.PHONY: test
test: build
	go test -v ./test

.PHONY: clean
clean:
	rm -r pwstore_backend

pwstore_backend: ${SOURCES}
	go build -o pwstore_backend cmd/rest/main.go

internal/api/api.gen.go: api/api.yaml configs/codegen/server.yaml
	mkdir -p internal/api
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config configs/codegen/server.yaml  api/api.yaml

internal/clients/google_open_id/google_open_id.gen.go: api/clients/google_open_id.yaml configs/codegen/google_open_id.yaml
	mkdir -p internal/clients/google_open_id
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config configs/codegen/google_open_id.yaml  api/clients/google_open_id.yaml
