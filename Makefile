PREFIX ?= /usr/local

.PHONY: build clean install

all: build

build:
	go build cmd/batteryhook.go

clean:
	rm batteryhook

install: build
	cp batteryhook ${PREFIX}/bin
