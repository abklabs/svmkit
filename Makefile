GO_DIRS		=	pkg cmd
GO_PACKAGES	=	$(foreach dir,$(GO_DIRS),./$(dir)/...)

.PHONY: lint test vet setup

test:
	go test $(GO_PACKAGES)

setup: 
	./setup

lint:
	./bin/check-env
	golangci-lint run $(GO_PACKAGES)
	shfmt -d .githooks/*
	shellcheck -P .githooks .githooks/*

vet:
	go vet $(GO_PACKAGES)

check: test lint

format:
	for dir in $(GO_DIRS) ; do ( cd $$dir && go fmt ./... ) ; done
	shfmt -w .githooks/*

clean:
	rm -f .env-checked	


.env-checked: bin/check-env
	./bin/check-env
	touch .env-checked

include .env-checked
