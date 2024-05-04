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