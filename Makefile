.PHONY: all build run clean

all: build

build:
	go build -o xkcd cmd/xkcd/main.go

run: build
	./xkcd

clean:
	rm -f xkcd

bench:
	go test -bench . -benchtime 10000x ./internal/search