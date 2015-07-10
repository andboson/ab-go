GOPATH=$(CURDIR)/../../

all:godeps  bindata.go build

godeps:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/gorilla/websocket

bindata.go: static/index.html static/main.js
	go-bindata -pkg='bindata' -o bindata/bindata.go static/

build:
	go build

clean:
	rm bindata/bindata.go