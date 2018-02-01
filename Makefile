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

release-linux: bin bindata.go
	go build -ldflags "-X main.Version=$(VERSION)" -o bin/pomo-$(VERSION)-linux-amd64

release-osx: bin bindata.go
	# This is used to cross-compile a Darwin compatible Mach-O executable 
	# on Linux for OSX, you need to install https://github.com/tpoechtrager/osxcross
	PATH="$$PATH:/usr/local/osx-ndk-x86/bin" GOOS=darwin GOARCH=amd64 CC=/usr/local/osx-ndk-x86/bin/x86_64-apple-darwin15-cc CGO_ENABLED=1 go build $(FLAGS) -o bin/pomo-$(VERSION)-darwin-amd64


docs: readme
	cd www && hugo -d ../docs

readme: www/data/readme.json

www/data/readme.json: www/data README.md
	cat README.md | python -c 'import json,sys; print(json.dumps({"content": sys.stdin.read()}))' > $@
www/data bin:
	mkdir -p $@
