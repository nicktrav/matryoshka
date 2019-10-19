#-----------------------------------------------------------------------------
# Target: setup
#-----------------------------------------------------------------------------

.PHONY: get
get: ## Fetch all dependencies
	go get ./...

#-----------------------------------------------------------------------------
# Target: tests
#-----------------------------------------------------------------------------

.PHONY: fmt
fmt: ## Format all files
	gofmt -w .

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: test
test: ## Run all tests
	go test -v ./...

.PHONY: test-race
test-race: ## Run the race detector
	go test -race -v ./...

test-all: fmt vet test test-race ## Run all tests, linters and formatters

#-----------------------------------------------------------------------------
# Target: artifacts
#-----------------------------------------------------------------------------

BINARY=matryoshka
BINDIR=$(CURDIR)/bin

bin: ## Creates the directory for binaries
	mkdir ${BINDIR}

VERSION=`git rev-parse HEAD`
LDFLAGS=-ldflags "-X github.com/nicktrav/matryoshka/cmd/version.BuildCommit=${VERSION}"

.PHONY: dist
dist: ${BINDIR} ## Builds the matryoshka binary
	go build ${LDFLAGS} -o ${BINDIR}/${BINARY} $(CURDIR)/cmd/main.go

.PHONY: install
install: dist ## Installs the binary into the GOBIN
	cp ${BINDIR}/${BINARY} ${GOBIN}/

#-----------------------------------------------------------------------------
# Target: help
#-----------------------------------------------------------------------------
.PHONY: help
help: ## Display this help text
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
