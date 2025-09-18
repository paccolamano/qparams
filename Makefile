PROJDIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

# change to project dir so we can express all as relative paths
$(shell cd $(PROJDIR))

VERSION ?= $(shell scripts/git-version.sh)

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: test
test:
	go test -v -cover -parallel 3 ./...

.PHONY: lint
lint:
	go tool golangci-lint run
