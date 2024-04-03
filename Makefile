.PHONY: build clean tool lint help

all: build

build:
	go build -o ./bin/a.out

clean:
	rm -rf bin

test_post:
	go test -v -run TestPost dcard-assignment/test/api/v1/ad

test_get:
	go test -v -run TestGet dcard-assignment/test/api/v1/ad