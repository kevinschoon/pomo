
.PHONY: \
	all 

all: bin/pomo

clean: 
	rm -v bin/pomo bindata.go

bindata.go:
	go-bindata -pkg main -o $@ tomato-icon.png

bin/pomo: bindata.go
	mkdir bin 2>/dev/null
	go build -o bin/pomo
