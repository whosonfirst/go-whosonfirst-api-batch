CWD=$(shell pwd)
GOPATH := $(CWD)

build:	rmdeps deps fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-api-batch; then rm -rf src/github.com/whosonfirst/go-whosonfirst-api-batch; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api-batch/http
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api-batch/lock
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api-batch/parse
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api-batch/process
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-api-batch/request
	cp *.go src/github.com/whosonfirst/go-whosonfirst-api-batch/
	cp http/*.go src/github.com/whosonfirst/go-whosonfirst-api-batch/http/
	cp lock/*.go src/github.com/whosonfirst/go-whosonfirst-api-batch/lock/
	cp parse/*.go src/github.com/whosonfirst/go-whosonfirst-api-batch/parse/
	cp process/*.go src/github.com/whosonfirst/go-whosonfirst-api-batch/process/
	cp request/*.go src/github.com/whosonfirst/go-whosonfirst-api-batch/request/
	if test ! -d src; then mkdir src; fi
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/patrickmn/go-cache"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-api"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-hash"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt http/*.go
	go fmt lock/*.go
	go fmt parse/*.go
	go fmt process/*.go
	go fmt request/*.go

bin:	self
	@GOPATH=$(shell pwd) go build -o bin/wof-api-batch-server cmd/wof-api-batch-server.go
