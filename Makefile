all: clean test

build:
	go mod tidy; go mod verify

test: build
	@echo "***** UNIT TESTS NOT YET PROVIDED *****"

clean:

.PHONY: all build test clean