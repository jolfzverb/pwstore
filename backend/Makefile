BUILD=build
SOURCES=$(shell find . -name "*.go")

.PHONY: format
format:
	gofmt -w ${SOURCES}

.PHONY: build
build: ${BUILD}/pwstore_backend

.PHONY: codegen
codegen: internal/api/api.go

.PHONY: run
run: build
	${BUILD}/pwstore_backend

.PHONY: test
test: build
	go test -v ./internal/tests

${BUILD}/pwstore_backend: ${SOURCES}
	mkdir -p ${BUILD}
	go build -o ${BUILD}/pwstore_backend cmd/rest/main.go

src/api/api.go: docs/api/api.yaml codegen/server.yaml
	mkdir -p internal/api
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config codegen/server.yaml  docs/api/api.yaml
