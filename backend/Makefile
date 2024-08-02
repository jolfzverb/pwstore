BUILD=build
SOURCES=$(shell find src/ -name "*.go")

.PHONY: format
format:
	gofmt -w ${SOURCES}

.PHONY: build
build: ${BUILD}/pwstore_backend

.PHONY: codegen
codegen: src/api/api.go

.PHONY: run
run: build
	./build/pwstore_backend

.PHONY: test
test: build
	go test -v ./src/tests

${BUILD}/pwstore_backend: ${SOURCES}
	mkdir -p ${BUILD}
	go build -o ${BUILD}/pwstore_backend src/main.go

src/api/api.go: docs/api/api.yaml codegen/server.yaml
	mkdir -p src/api
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config codegen/server.yaml  docs/api/api.yaml
