DOCKER_CMD=docker run --rm -ti --user 1000 -w /build/pomo -v $$PWD:/build/pomo
DOCKER_IMAGE=pomo-build
VERSION ?= $(shell git describe --tags 2>/dev/null)
ifeq "$(VERSION)" ""
	VERSION := UNKNOWN
endif

.PHONY: \
	test \
	docs \
	pomo-build \
	readme \
	release \
	release-linux \
	release-darwin

bin/pomo: test
	go build -o $@

bindata.go: tomato-icon.png
	go-bindata -pkg main -o $@ $^

test:
	go test ./...
	go vet ./...

pomo-build:
	docker build -t $(DOCKER_IMAGE) .

bin/pomo-linux: bin/pomo-$(VERSION)-linux-amd64

bin/pomo-darwin: bin/pomo-$(VERSION)-darwin-amd64

bin/pomo-$(VERSION)-linux-amd64: bin bindata.go
	$(DOCKER_CMD) --env GOOS=linux --env GOARCH=amd64 $(DOCKER_IMAGE) go build -ldflags "-X main.Version=$(VERSION)" -o $@

bin/pomo-$(VERSION)-linux-amd64.md5:
	md5sum bin/pomo-$(VERSION)-linux-amd64 | sed -e 's/bin\///' > $@

bin/pomo-$(VERSION)-darwin-amd64: bin bindata.go
	# This is used to cross-compile a Darwin compatible Mach-O executable
	# on Linux for OSX, you need to install https://github.com/tpoechtrager/osxcross
	$(DOCKER_CMD) --env GOOS=darwin --env GOARCH=amd64 --env CC=x86_64-apple-darwin15-cc --env CGO_ENABLED=1 $(DOCKER_IMAGE) go build -ldflags "-X main.Version=$(VERSION)" -o $@


bin/pomo-$(VERSION)-darwin-amd64.md5:
	md5sum bin/pomo-$(VERSION)-darwin-amd64 | sed -e 's/bin\///' > $@

release-linux: bin/pomo-$(VERSION)-linux-amd64 bin/pomo-$(VERSION)-linux-amd64.md5

release-darwin: bin/pomo-$(VERSION)-darwin-amd64 bin/pomo-$(VERSION)-darwin-amd64.md5

release: release-linux release-darwin

docs: www/data/readme.json
	cd www && cp ../install.sh static/ && hugo -d ../docs

www/data/readme.json: www/data README.md
	cat README.md | python -c 'import json,sys; print(json.dumps({"content": sys.stdin.read()}))' > $@

www/data bin:
	mkdir -p $@
