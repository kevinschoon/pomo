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

bindata.go: tomato-icon.png
	go-bindata -pkg main -o $@ $^

test:
	go test ./...
	go vet ./...

release: bin bindata.go
	go build -ldflags "-X main.Version=$(VERSION)" -o bin/pomo-$(VERSION)-linux

docs: readme
	cd www && hugo -d ../docs

readme: www/data/readme.json

www/data/readme.json: www/data README.md
	cat README.md | python -c 'import json,sys; print(json.dumps({"content": sys.stdin.read()}))' > $@
www/data bin:
	mkdir -p $@
