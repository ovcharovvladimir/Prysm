#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
if [[ $PWD = *"travis"* ]];
then
    solution="$GOPATH/go/src"
    #"${PWD%/go/src/*}"
else
    solution="${PWD%/go/src/*}/go/src"
fi

#"$GOPATH"
root="$PWD"
echo "GOPATH: $GOPATH"
echo "Project Working Dir: $PWD"
echo "Workspace: $workspace"
echo "Solution: $solution"
dir="$workspace/src/github.com/ovcharovvladimir"

if [ ! -L "$dir/Prysm" ]; then
    mkdir -p "$dir"
    cd "$dir"
    ln -s ../../../../../. Prysm
    cd "$root"
fi
if [ ! -L "$dir/essentiaHybrid" ]; then
    mkdir -p "$dir"
    cd "$dir"
    ln -s $solution/github.com/ovcharovvladimir/essentiaHybrid/. essentiaHybrid
    cd "$root"
fi
pth="$workspace/src/google.golang.org"

if [ ! -L "$pth/grpc" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/google.golang.org/grpc/. grpc
    cd "$root"
fi
pth="$workspace/src/google.golang.org"

if [ ! -L "$pth/genproto" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/google.golang.org/genproto/. genproto
    cd "$root"
fi
#golang.org/x/sys/unix
pth="$workspace/src/golang.org"

if [ ! -L "$pth/x" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/golang.org/x/. x
    cd "$root"
fi
#github.com/sirupsen/logrus
pth="$workspace/src/github.com/"

if [ ! -L "$pth/sirupsen" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/sirupsen/. sirupsen
    cd "$root"
fi
#
#github.com/x-cray/logrus-prefixed-formatter
pth="$workspace/src/github.com"

if [ ! -L "$pth/x-cray" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/x-cray/. x-cray
    cd "$root"
fi
#
#github.com/urfave/cli
pth="$workspace/src/github.com"

if [ ! -L "$pth/urfave" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/urfave/. urfave
    cd "$root"
fi
#
#github.com/syndtr/goleveldb/leveldb/errors
pth="$workspace/src/github.com"

if [ ! -L "$pth/syndtr" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/syndtr/. syndtr
    cd "$root"
fi
#
#github.com/multiformats/go-multiaddr
pth="$workspace/src/github.com"

if [ ! -L "$pth/multiformats" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/multiformats/. multiformats
    cd "$root"
fi
#
#github.com/libp2p/go-libp2p/p2p/discovery
pth="$workspace/src/github.com"

if [ ! -L "$pth/libp2p" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/libp2p/. libp2p
    cd "$root"
fi
#
#github.com/whyrusleeping/timecache
pth="$workspace/src/github.com"

if [ ! -L "$pth/whyrusleeping" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/whyrusleeping/. whyrusleeping
    cd "$root"
fi
#
#github.com/spaolacci/murmur3
pth="$workspace/src/github.com"

if [ ! -L "$pth/spaolacci" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/spaolacci/. spaolacci
    cd "$root"
fi
#
#github.com/btcsuite/btcd/btcec
pth="$workspace/src/github.com"

if [ ! -L "$pth/btcsuite" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/btcsuite/. btcsuite
    cd "$root"
fi
#
#github.com/coreos/go-semver/semver
pth="$workspace/src/github.com"

if [ ! -L "$pth/coreos" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/coreos/. coreos
    cd "$root"
fi
#
#github.com/minio/sha256-simd
pth="$workspace/src/github.com"

if [ ! -L "$pth/minio" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/minio/. minio
    cd "$root"
fi
#
#github.com/fjl/memsize/memsizeui
pth="$workspace/src/github.com"

if [ ! -L "$pth/fjl" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/fjl/. fjl
    cd "$root"
fi
#
#github.com/jbenet
pth="$workspace/src/github.com"

if [ ! -L "$pth/jbenet" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/jbenet/. jbenet
    cd "$root"
fi
#
#github.com/ipfs
pth="$workspace/src/github.com"

if [ ! -L "$pth/ipfs" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/ipfs/. ipfs
    cd "$root"
fi
#
#github.com/opentracing
pth="$workspace/src/github.com"

if [ ! -L "$pth/opentracing" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/opentracing/. opentracing
    cd "$root"
fi
#
#github.com/golang
pth="$workspace/src/github.com"

if [ ! -L "$pth/golang" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/golang/. golang
    cd "$root"
fi
#
#github.com/fd
pth="$workspace/src/github.com"

if [ ! -L "$pth/fd" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/fd/. fd
    cd "$root"
fi
#
#github.com/gogo
pth="$workspace/src/github.com"

if [ ! -L "$pth/gogo" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/gogo/. gogo
    cd "$root"
fi
#
#github.com/google
pth="$workspace/src/github.com"

if [ ! -L "$pth/google" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/google/. google
    cd "$root"
fi
#
#github.com/gorilla
pth="$workspace/src/github.com"

if [ ! -L "$pth/gorilla" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/gorilla/. gorilla
    cd "$root"
fi
#
#github.com/gxed
pth="$workspace/src/github.com"

if [ ! -L "$pth/gxed" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/gxed/. gxed
    cd "$root"
fi
#
#github.com/huin
pth="$workspace/src/github.com"

if [ ! -L "$pth/huin" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/huin/. huin
    cd "$root"
fi
#
#github.com/jackpal
pth="$workspace/src/github.com"

if [ ! -L "$pth/jackpal" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/jackpal/. jackpal
    cd "$root"
fi
#
#github.com/mattn
pth="$workspace/src/github.com"

if [ ! -L "$pth/mattn" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/mattn/. mattn
    cd "$root"
fi
#
#
#github.com/mgutz
pth="$workspace/src/github.com"

if [ ! -L "$pth/mgutz" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/mgutz/. mgutz
    cd "$root"
fi
#
#github.com/miekg
pth="$workspace/src/github.com"

if [ ! -L "$pth/miekg" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/miekg/. miekg
    cd "$root"
fi
#
#github.com/mr-tron
pth="$workspace/src/github.com"

if [ ! -L "$pth/mr-tron" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/mr-tron/. mr-tron
    cd "$root"
fi
#
echo "Path $dir"
# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$dir/Prysm"
PWD="$dir/Prysm"
echo "----- $@"
# Launch the arguments with the configured environment.
exec "$@"
