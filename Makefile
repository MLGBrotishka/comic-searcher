.PHONY: all build run clean

all: build

run: build
	./xkcd

build:
	go build -o xkcd cmd/xkcd/main.go

clean:
	rm -f xkcd

bench:
	go test -bench . -benchtime 10000x ./internal/search