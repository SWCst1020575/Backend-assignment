.PHONY: build clean tool lint help

all: build

build:
	go build -o ./bin/a.out

clean:
	rm -rf bin