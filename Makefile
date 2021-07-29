DOCKER_CMD=docker run --rm -ti -w /build/pomo -v $$PWD:/build/pomo
DOCKER_IMAGE=pomo-build

VERSION ?= $(shell git describe --tags 2>/dev/null)
ifeq "$(VERSION)" ""
	VERSION := UNKNOWN
endif

LDFLAGS=\
	-X github.com/kevinschoon/pomo/pkg/internal/version.Version=$(VERSION)

.PHONY: \
	test \
	docs \
	pomo-build \
	readme 

default:
	cd cmd/pomo && \
	go install -ldflags '${LDFLAGS}'

bin/pomo: test
	cd cmd/pomo && \
	go build -ldflags '${LDFLAGS}' -o ../../$@

#bindata.go: tomato-icon.png
#	go-bindata -pkg main -o $@ $^

test:
	go test ./...
	go vet ./...

docs: www/data/readme.json
	cd www && hugo -d ../docs

www/data/readme.json: www/data README.md
	cat README.md | python -c 'import json,sys; print(json.dumps({"content": sys.stdin.read()}))' > $@

www/data bin:
	mkdir -p $@
