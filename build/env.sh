#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
gopath= "$GOPATH"
root="$PWD"
essdir="$gopath/src/github.com/ovcharovvladimir"
dir="$workspace/src/github.com/ovcharovvladimir"

if [ ! -L "$dir/Prysm" ]; then
    mkdir -p "$dir"
    cd "$dir"
    ln -s ../../../../../. Prysm
    cd "$root"
fi
if [ ! -L "$essdir/essentiaHybrid" ]; then
    mkdir -p "$essdir"
    cd "$essdir"
    ln -s ../../../../../. essentiaHybrid
    cd "$root"
fi
# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$dir/Prysm"
PWD="$dir/Prysm"

# Launch the arguments with the configured environment.
exec "$@"
