.PHONY: lint test vet

test:
	go test ./pkg/...

lint:
	golangci-lint run ./pkg/...

vet:
	go vet ./pkg/...
