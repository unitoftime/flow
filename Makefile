all: test

generate:
	go generate ./...

fmt:
	go fmt ./...

test:
	go test ./...

upgrade:
	go get -u github.com/unitoftime/ecs@HEAD
	go get -u ./...
	go mod tidy

coverage: fmt
	go test ./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
