#!/usr/bin/env bash

set -euo pipefail

readonly src_dir=/tmp/reserve-stats-abi-gen/src/github.com/ethereum
export GOPATH=/tmp/reserve-stats-abi-gen
export PATH=$GOPATH/bin:$PATH

mkdir -p "$src_dir"
cd "$src_dir"
[[ -d go-ethereum ]] || git clone https://github.com/ethereum/go-ethereum.git
go install -v github.com/ethereum/go-ethereum/cmd/abigen


abigen -abi "$OLDPWD"/internal_network.abi -pkg contracts -type InternalNetwork -out "$OLDPWD"/internal_network_abi.go
abigen -abi "$OLDPWD"/wrapper.abi -pkg contracts -type Wrapper -out "$OLDPWD"/wrapper_abi.go
abigen -abi "$OLDPWD"/reserve.abi -pkg contracts -type Reserve -out "$OLDPWD"/reserve.go
abigen -abi "$OLDPWD"/sanity_rates.abi -pkg contracts -type SanityRates -out "$OLDPWD"/sanity_rates.go
abigen -abi "$OLDPWD"/conversion_rates.abi -pkg contracts -type ConversionRates -out "$OLDPWD"/conversion_rates.go
abigen -abi "$OLDPWD"/erc20.abi -pkg contracts -type ERC20 -out "$OLDPWD"/erc20.go # erc20 type 1 where symbol return in string format
abigen -abi "$OLDPWD"/erc20_bytes32.abi -pkg contracts -type ERC20Type2 -out "$OLDPWD"/erc20_type2.go # erc20 type2 where symbol return in bytes32 format

