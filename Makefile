
test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out >> /dev/null
	go tool cover -func coverage.out

lint:
	golangci-lint run --enable-all