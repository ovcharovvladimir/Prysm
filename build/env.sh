#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
solution="${PWD%/go/src/*}/go/src/github.com/ovcharovvladimir"
#"$GOPATH"
root="$PWD"

echo "Workspace $workspace"
echo "Solution $solution"
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
    ln -s $solution/essentiaHybrid/. essentiaHybrid
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
