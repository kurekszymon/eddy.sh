build:
	go mod tidy
	go build

clean: 
	go clean

test: 
	go test ./... -timeout 90s

lint:
	golangci-lint run