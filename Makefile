# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: android ios sness  all  clean
#.PHONY: sn-linux sn-linux-386 sn-linux-amd64 sn-linux-mips64 sn-linux-mips64le
#.PHONY: sn-linux-arm sn-linux-arm-5 sn-linux-arm-6 sn-linux-arm-7 sn-linux-arm64
#.PHONY: sn-darwin sn-darwin-386 sn-darwin-amd64
.PHONY: sn-windows sn-windows-386 sn-windows-amd64

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
sn-cross: sness-linux sness-darwin sness-windows #sness-android sness-ios
	@echo "Full cross compilation done"

sn-linux: sn-linux-386 sn-linux-amd64 sn-linux-arm #sn-linux-mips sn-linux-mips64 sn-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/linux/386/sness
	@ls -ld $(GOBIN)/linux/amd64/sness
	@ls -ld $(GOBIN)/linux/arm5/sness
	@ls -ld $(GOBIN)/linux/arm6/sness
	@ls -ld $(GOBIN)/linux/arm7/sness
	#@ls -ld (GOBIN)/linux/mips/sness-linux-mips/sness
	#@ls -ld (GOBIN)/linux/mips/sness-linux-mipsle/sness




sn-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/386 --targets=linux/386 -v ./sness 
	@echo "Linux 386 cross compilation done:"
	@sudo mv $(GOBIN)/linux/386/sness-linux-386 $(GOBIN)/linux/386/sness
	@sudo rm -fv  $(GOBIN)/linux/386/sness-linux-386

sn-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/amd64 -targets=linux/amd64 -v ./sness 
	@echo "Linux amd64 cross compilation done:"
	@sudo mv $(GOBIN)/linux/amd64/sness-linux-amd64 $(GOBIN)/linux/amd64/sness
	@sudo rm -fv  $(GOBIN)/linux/amd64/sness-linux-amd64

sn-linux-arm: sn-linux-arm-5 sn-linux-arm-6 sn-linux-arm-7 sn-linux-arm64
	@echo "Linux ARM cross compilation done"

sn-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm5 --targets=linux/arm-5 -v ./sness
	@echo "Linux ARMv5 cross compilation done:"
	@sudo mv $(GOBIN)/linux/arm5/sness-linux-arm-5 $(GOBIN)/linux/arm5/sness
	@sudo rm -fv  $(GOBIN)/linux/arm5/sness-linux-arm-5

sn-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm6 --targets=linux/arm-6 -v ./sness
	@echo "Linux ARMv6 cross compilation done:"
	@sudo mv $(GOBIN)/linux/arm6/sness-linux-arm-6 $(GOBIN)/linux/arm6/sness
	@sudo rm -fv  $(GOBIN)/linux/arm6/sness-linux-arm-6

sn-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm7 --targets=linux/arm-7 -v ./sness
	@echo "Linux ARMv7 cross compilation done:"
	@sudo mv $(GOBIN)/linux/arm7/sness-linux-arm-7 $(GOBIN)/linux/arm7/sness
	@sudo rm -fv  $(GOBIN)/linux/arm7/sness-linux-arm-7

sn-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/arm64 --targets=linux/arm64 -v ./sness
	@echo "Linux ARM64 cross compilation done:"
	@sudo mv $(GOBIN)/linux/arm64/sness-linux-arm64 $(GOBIN)/linux/arm64/sness
	@sudo rm -fv  $(GOBIN)/linux/arm64/sness-linux-arm64

sn-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/mips --targets=linux/mips --ldflags '-extldflags "-static"' -v ./sness
	@echo "Linux MIPS cross compilation done:"
	@sudo mv $(GOBIN)/linux/mips/sness-linux-mips $(GOBIN)/linux/mips/sness
	@sudo rm -fv  $(GOBIN)/linux/mips/sness-linux-mips

sn-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/mipsle  --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./sness
	@echo "Linux MIPSle cross compilation done:"
	@sudo mv $(GOBIN)/linux/mipsle/sness-linux-mipsle $(GOBIN)/linux/mipsle/sness
	@sudo rm -fv  $(GOBIN)/linux/mipsle/sness-linux-mipsle

sn-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/linux/mips64   --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./sness
	@echo "Linux MIPS64 cross compilation done:"
	@sudo mv $(GOBIN)/linux/mips64/sness-linux-mips64 $(GOBIN)/linux/mips64/sness
	@sudo rm -fv  $(GOBIN)/linux/mips64/sness-linux-mips64

sn-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO)  --dest=$(GOBIN)/linux/mips64le  --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./sness
	@echo "Linux MIPS64le cross compilation done:"
	@sudo mv $(GOBIN)/linux/mips64le/sness-linux-mips64le $(GOBIN)/linux/mips64le/sness
	@sudo rm -fv  $(GOBIN)/linux/mips64le/sness-linux-mips64le
	

sn-darwin: sn-darwin-386 sn-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld  $(GOBIN)/darwin/386/sness
	@ls -ld  $(GOBIN)/darwin/amd64/sness

sn-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/darwin/386 --targets=darwin/386 -v ./sness
	@echo "Darwin 386 cross compilation done:"
	@sudo mv $(GOBIN)/darwin/386/sness-darwin-10.6-386 $(GOBIN)/darwin/386/sness
	@sudo rm -fv  $(GOBIN)/darwin/386/sness-darwin-10.6-386

sn-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/darwin/amd64 --targets=darwin/amd64 -v ./sness
	@echo "Darwin amd64 cross compilation done:"
	@sudo mv $(GOBIN)/darwin/amd64/sness-darwin-10.6-amd64 $(GOBIN)/darwin/amd64/sness
	@sudo rm -fv  $(GOBIN)/darwin/amd64/sness-darwin-10.6-amd64
	
sn-windows: sn-windows-386 sn-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld  $(GOBIN)/windows/386/sness.exe
	@ls -ld  $(GOBIN)/windows/amd64/sness.exe

sn-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/windows/386 --targets=windows/386 -v ./sness
	@echo "Windows 386 cross compilation done:"
	@sudo mv $(GOBIN)/windows/386/sness-windows-4.0-386.exe $(GOBIN)/windows/386/sness.exe
	@sudo rm -fv  $(GOBIN)/windows/386/sness-windows-4.0-386.exe

sn-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN)/windows/amd64 --targets=windows/amd64 -v ./sness
	@echo "Windows amd64 cross compilation done:"
	@sudo mv $(GOBIN)/windows/amd64/sness-windows-4.0-amd64.exe $(GOBIN)/windows/amd64/sness.exe
	@sudo rm -fv  $(GOBIN)/windows/amd64/sness-windows-4.0-amd64.exe
