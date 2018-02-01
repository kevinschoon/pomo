VERSION ?= $(shell git describe --tags 2>/dev/null)
ifeq "$(VERSION)" ""
	VERSION := UNKNOWN
endif

.PHONY: \
	all \
	test \
	docs \
	readme \
	release

all: bin/pomo

clean: 
	-rm -fv bin/* docs/*

bindata.go:
	go-bindata -pkg main -o $@ tomato-icon.png

test:
	go test ./...
	go vet ./...

release: bindata.go
	@echo mkdir bin 2>/dev/null || true
	go build -ldflags "-X main.Version=$(VERSION)" -o bin/pomo-$(VERSION)-linux

docs: readme
	cd www && hugo -d ../docs

readme: www/data
	cat README.md | python -c 'import json,sys; print(json.dumps({"content": sys.stdin.read()}))' > www/data/readme.json

www/data:
	mkdir -p $@
