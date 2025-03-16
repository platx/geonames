lint:
	golangci-lint --timeout=3m run -v

test:
	go test --tags=unit ./...

coverage:
	go test --tags=unit -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
