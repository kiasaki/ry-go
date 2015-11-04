deps:
	go get github.com/jteeuwen/go-bindata/...


build:
	#go-bindata -pkg ry lisp/
	go build -o ry cmd/ry/main.go

build-ry-repl:
	#go-bindata -pkg ry lisp/
	go build -o ry-repl cmd/ry-repl/main.go

test: build-ry-repl
	tests/testall.sh

setup:
	go get github.com/jteeuwen/go-bindata/go-bindata

clean:
	rm gobindata.go
