.PHONY: run build test clean

run:
	go run ./cmd/timer

build:
	go build -o timer ./cmd/timer

test:
	go test ./...

clean:
	rm -f timer
