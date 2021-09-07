build:
	go build -o target/tmihooks .
lint:
	golangci-lint run . internal/...
test:
	go test ./...