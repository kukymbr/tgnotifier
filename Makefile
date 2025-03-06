GOLANGCI_LINT_VERSION := 1.64.6

all:
	make clean
	make generate
	make tidy
	make test
	make build
	make lint

generate:
	go generate ./cmd/tgnotifier

install_tools:
	go install ./...

apis:
	protoc -I api/grpc api/grpc/tgnotifier.proto --go_out=./internal/api/grpc/ --go_opt=paths=source_relative --go-grpc_out=./internal/api/grpc/ --go-grpc_opt=paths=source_relative --oas_out=./api/http/
	oapi-codegen --config=./api/http/oapi.config.yml ./api/http/openapi.yaml

tidy:
	go mod tidy
	go vet ./...

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

build:
	go build ./cmd/tgnotifier

clean:
	go clean