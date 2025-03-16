lint:
	golangci-lint --timeout=3m run -v

test:
	go test -v -race ./...

cover:
	go test -coverprofile=coverage.txt
