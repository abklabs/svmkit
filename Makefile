GO_DIRS		=	pkg cmd
GO_PACKAGES	=	$(foreach dir,$(GO_DIRS),./$(dir)/...)

.PHONY: lint test vet setup

test:
	go test $(GO_PACKAGES)

setup: 
	./setup

lint:
	golangci-lint run $(GO_PACKAGES)

vet:
	go vet $(GO_PACKAGES)

check: test lint

format:
	for dir in $(GO_DIRS) ; do ( cd $$dir && go fmt ./... ) ; done
