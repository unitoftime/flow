all: test

fmt:
	go fmt ./...

test: fmt
	go test ./...

upgrade:
	go get -u ./...
	go mod tidy

coverage: fmt
	go test ./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
