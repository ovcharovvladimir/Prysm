# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: android ios sness  all  clean
#.PHONY: gess-linux gess-linux-386 gess-linux-amd64 gess-linux-mips64 gess-linux-mips64le
#.PHONY: gess-linux-arm gess-linux-arm-5 gess-linux-arm-6 gess-linux-arm-7 gess-linux-arm64
#.PHONY: gess-darwin gess-darwin-386 gess-darwin-amd64
#.PHONY: gess-windows gess-windows-386 gess-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

sness:
	bash build/env.sh go run build/ci.go install ./sness
	@echo "Done building Ess Supernode."
	@echo "Run \"$(GOBIN)/sness\" to launch sness"

all:
	bash build/env.sh go run build/ci.go install
	@echo "Run \"$(GOBIN)"


clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u gopkg.in/urfave/cli.v1
	env GOBIN= go get -u github.com/robertkrimen/otto
	env GOBIN= go get -u github.com/fatih/color
	env GOBIN= go get -u github.com/docker/docker/pkg/reexec
	env GOBIN= go get -u github.com/Azure/azure-storage-blob-go/2018-03-28/azblob
	env GOBIN= go get -u github.com/x-cray/logrus-prefixed-formatter
	env GOBIN= go get -u  github.com/multiformats/go-multiaddr
	env GOBIN= go get -u  github.com/urfave/cli
	env GOBIN= go get -u google.golang.org/grpc
	env GOBIN= go get -u golang.org/x/crypto/blake2b
	env GOBIN= go get -u github.com/libp2p/go-libp2p/
	env GOBIN= go get -u github.com/fjl/memsize/memsizeui
	env GOBIN= go get -u github.com/libp2p/go-floodsub
	env GOBIN= go get -u github.com/whyrusleeping/mdns
	env GOBIN= go get -u github.com/syndtr/goleveldb/leveldb
	env GOBIN= go get -u -u github.com/golang/protobuf/proto


# Cross Compilation Targets (xgo)
