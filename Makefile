.PHONY: all
all: build

.PHONY: build
build:
	mkdir build || true
	go build -o build/ ./client/*.go
	go build -o build/ ./server/*.go

.PHONY: clean
clean:
	rm -r build/

.PHONY: test
test:
	@true
