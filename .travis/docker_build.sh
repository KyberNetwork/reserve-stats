#!/bin/bash
# -*- firestarter: "shfmt -i 4 -ci -w %p" -*-

set -euxo pipefail

readonly build_part=${BUILD_PART:-}
readonly golangci_config=$(readlink -f .golangci.yml)

build() {
    local build_dir="$1"
    pushd "$build_dir"
    golangci-lint run --config ${golangci_config} -v
    go test -v -race -mod=vendor ./...
    popd
}

# build_file loads and builds the configuration from given file
build_file() {
    local config_file="$1"
    while read -r line; do
        build $line
    done < "$config_file"
}


build_docs() {
    pushd ./apidocs
    docker-compose up
    popd
}

case "$build_part" in
    1)
        build_file .travis/build_part_1
        ;;
    2)
        build_file .travis/build_part_2
        ;;
    *)
        # find list of modules that already build in above step
        modules=($(sed -e 's/ .*$//' .travis/build_part_*))
        # join module arrays with \|
        exclude_pattern=$(printf '\|%s' "${modules[@]}")
        # remove leading \|
        exclude_pattern=${exclude_pattern:2}
        exclude_pattern=$(printf 'github.com/KyberNetwork/reserve-stats/\(%s\)' "$exclude_pattern")
        golangci-lint run --config ${golangci_config}  --exclude "$exclude_pattern"
        go test -v -race -mod=vendor $(go list -mod=vendor ./... | grep -v "$exclude_pattern")

        if [[ $TRAVIS_BRANCH == 'develop' ]]; then
            build_docs
        fi
        ;;
esac
