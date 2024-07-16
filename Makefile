all:
	make clean
	make generate
	make tidy
	make test
	make build
	make lint

generate:
	go generate ./cmd/tgnotifier

tidy:
	go mod tidy
	go vet ./...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
	golangci-lint run ./...

test:
	go test -race -coverprofile=coverage_out ./...
	go tool cover -func=coverage_out
	go tool cover -html=coverage_out -o coverage.html
	rm -f coverage_out

test_short:
	go test -short ./...

build:
	go build ./cmd/tgnotifier

clean:
	go clean