include .env.example
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

.PHONY: all build run clean

all: run

run: ## tidy run
	go mod tidy && go mod download && \
	CGO_ENABLED=1 go run -tags migrate ./cmd/app
.PHONY: run

build: ## build
	go build -o $(LOCAL_BIN)/xkcd cmd/app/main.go
.PHONY: build

clean: ## removes local build
	rm -f $(LOCAL_BIN)/xkcd 
.PHONY: clean

migrate-create:  ## create new migration
	migrate create -ext sql -dir migrations 'migrate_name'
.PHONY: migrate-create

migrate-up: ## migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

sec: build ## check vulnerability
	go install golang.org/x/vuln/cmd/govulncheck@latest && \
	govulncheck ./...
.PHONY: sec

test:
	go test ./... -coverprofile=coverage.out -race && go tool cover -html=coverage.out
.PHONY: test

e2e:
	./e2e/test.sh
.PHONY: e2e

# Linter
GOLANGCI_LINT = golangci-lint  # location $(LOCAL_BIN)/golangci-lint 

.install-linter:
	$(shell [ -f bin ] || mkdir -p $(LOCAL_BIN))
	### INSTALL GOLANGCI-LINT ###
	[ -f $(GOLANGCI_LINT) ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN) v1.57.2
.PHONY: .install-linter

lint:
	### RUN GOLANGCI-LINT ###
	$(GOLANGCI_LINT) run ./... --config=.golangci.yaml
.PHONY: lint

lint-fast:
	$(GOLANGCI_LINT) run ./... --fast --config=.golangci.yaml
.PHONY: lint-fast