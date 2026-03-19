BINARY    := git-profile
VERSION   := $(shell cat VERSION)
LDFLAGS   := -ldflags "-s -w -X main.version=$(VERSION)"

# Capture everything after `make run` as arguments
# Usage: make run show / make run list / make run use personal --global
ifeq (run,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif

.PHONY: build install local-install clean fmt vet test run help

build: ## Build the binary
	go build $(LDFLAGS) -o $(BINARY) .

install: ## Install to $GOPATH/bin
	go install $(LDFLAGS) .

local-install: build ## Build and install to /usr/local/bin
	sudo cp $(BINARY) /usr/local/bin/$(BINARY)

clean: ## Remove built binary
	rm -f $(BINARY)

fmt: ## Format source code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

test: ## Run tests
	go test ./...

run: build ## Build and run (e.g. make run show, make run list)
	./$(BINARY) $(RUN_ARGS)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
