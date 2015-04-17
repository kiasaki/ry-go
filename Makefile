build-ry:
	go-bindata -pkg ry lisp/
	go build -o ry cmd/ry/main.go

build-ryl:
	go-bindata -pkg ry lisp/
	go build -o ryl cmd/ryl/main.go

test:
	tests/testall.sh

setup:
	go get github.com/jteeuwen/go-bindata/go-bindata

clean:
	rm gobindata.go
