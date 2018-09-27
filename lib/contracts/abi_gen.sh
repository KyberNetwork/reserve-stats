#!/usr/bin/env bash

set -euox pipefail

export GOPATH=/tmp/reserve-stats-abi-gen
export PATH=$GOPATH/bin:$PATH

mkdir -p /tmp/reserve-stats-abi-gen/src/github.com/ethereum
cd /tmp/reserve-stats-abi-gen/src/github.com/ethereum
[[ -d go-ethereum ]] || git clone https://github.com/ethereum/go-ethereum.git
go install -v github.com/ethereum/go-ethereum/cmd/abigen

abigen -abi "$OLDPWD"/internal_network.abi -pkg contracts -type InternalNetwork -out "$OLDPWD"/internal_network_abi.go
