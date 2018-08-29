#!/usr/bin/env bash

declare -a skip_coverage=("github.comgithub.com/ovcharovvladimir/Prysm/contracts/sharding-manager-contract"
                          "github.com/github.com/ovcharovvladimir/Prysm/voter-registration-contract")

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    if [[ ${skip_coverage[*]} =~ "$d" ]]; then
        continue
    fi
    go test -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
