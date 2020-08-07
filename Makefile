PREFIX = /usr/local

all: build

build:
	go build cmd/batteryhook.go

clean:
	rm batteryhook

install: build
	install -Dm0755 batteryhook ${PREFIX}/bin/batteryhook

.PHONY: build clean install
