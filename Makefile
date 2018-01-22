VERSION ?= $(shell git describe --tags 2>/dev/null)
ifeq "$(VERSION)" ""
	VERSION := UNKNOWN
endif

.PHONY: \
	all \
	test

all: bin/pomo

clean: 
	rm -v bin/* 2> /dev/null || true

bindata.go:
	go-bindata -pkg main -o $@ tomato-icon.png

test:
	go test ./...

bin/pomo: bindata.go
	@echo mkdir bin 2>/dev/null || true
	go build -ldflags "-X main.Version=$(VERSION)" -o bin/pomo

www/data/readme.json:
	cat README.md | python -c 'import json,sys; print(json.dumps({"content": sys.stdin.read()}))' > $@
