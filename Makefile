GOLANGCI_LINT_VERSION := 1.64.6

all:
	$(MAKE) clean
	$(MAKE) prepare
	$(MAKE) validate
	$(MAKE) build
	$(MAKE) build_gui

prepare:
	go mod tidy
	go install ./...
	$(MAKE) apis

validate:
	go vet ./...
	$(MAKE) lint
	$(MAKE) test

build:
	go build $(build_arguments) ./cmd/tgnotifier

build_without_gprc:
	go build -tags no_grpc ./cmd/tgnotifier

build_without_http:
	go build -tags no_http ./cmd/tgnotifier

build_without_servers:
	go build -tags no_http,no_grpc ./cmd/tgnotifier

build_gui:
	go build -tags gui,no_http,no_grpc ./cmd/tgnotifierui

apis:
	protoc -I api/grpc api/grpc/tgnotifier.proto --go_out=./internal/api/grpc/ --go_opt=paths=source_relative --go-grpc_out=./internal/api/grpc/ --go-grpc_opt=paths=source_relative --oas_out=./api/http/
	oapi-codegen --config=./api/http/oapi.config.yml ./api/http/openapi.yaml

lint:
	if [ ! -f ./bin/golangci-lint ]; then \
  		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "./bin" "v${GOLANGCI_LINT_VERSION}"; \
  	fi;
	./bin/golangci-lint run ./...

test:
	go test -race -coverprofile=coverage_out ./...
	go tool cover -func=coverage_out
	go tool cover -html=coverage_out -o coverage.html
	rm -f coverage_out

test_short:
	go test -short ./...

test_grpc:
	go test ./internal/api/tests/... -v -tags grpc_tests -count 1

clean:
	go clean