.DELETE_ON_ERROR:

# High-level commands
.PHONY: build
build: bin bin/djafka

bin:
	mkdir -p bin

bin/djafka: .FORCE
	go build -o bin/djafka

.PHONY: download
download:
	go mod download

.PHONY: clean
clean:
	rm -f bin/*

.FORCE:
