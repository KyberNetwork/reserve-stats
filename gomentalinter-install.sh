#!/bin/bash
 FILE_LIST=(\
     "gometalinter-2.0.11-linux-amd64/gometalinter" \
     "gometalinter-2.0.11-linux-amd64/gocyclo" \
     "gometalinter-2.0.11-linux-amd64/nakedret" \
     "gometalinter-2.0.11-linux-amd64/misspell" \
     "gometalinter-2.0.11-linux-amd64/gosec" \
     "gometalinter-2.0.11-linux-amd64/golint" \
     "gometalinter-2.0.11-linux-amd64/ineffassign" \
     "gometalinter-2.0.11-linux-amd64/goconst" \
     "gometalinter-2.0.11-linux-amd64/errcheck" \
     "gometalinter-2.0.11-linux-amd64/maligned" \
     "gometalinter-2.0.11-linux-amd64/unconvert" \
     "gometalinter-2.0.11-linux-amd64/dupl" \
     "gometalinter-2.0.11-linux-amd64/structcheck" \
     "gometalinter-2.0.11-linux-amd64/varcheck" \
     "gometalinter-2.0.11-linux-amd64/safesql" \
     "gometalinter-2.0.11-linux-amd64/deadcode" \
     "gometalinter-2.0.11-linux-amd64/lll" \
     "gometalinter-2.0.11-linux-amd64/goimports" \
     "gometalinter-2.0.11-linux-amd64/gotype" \
     "gometalinter-2.0.11-linux-amd64/gosimple" \
     "gometalinter-2.0.11-linux-amd64/megacheck" \
     "gometalinter-2.0.11-linux-amd64/staticcheck" \
     "gometalinter-2.0.11-linux-amd64/unused" \
     "gometalinter-2.0.11-linux-amd64/interfacer" \
     "gometalinter-2.0.11-linux-amd64/unparam" \
     "gometalinter-2.0.11-linux-amd64/gochecknoinits" \
     "gometalinter-2.0.11-linux-amd64/gochecknoglobals" \
)
CURRENT_DIR=$PWD
mkdir -p ${TRAVIS_HOME}/gometalinter
cd ${TRAVIS_HOME}/gometalinter
for file in $FILE_LIST; do
    if ! [ -f $file ]; then
        rm -rf gometalinter-2.0.11-linux-amd64*
        curl -O -J -L https://github.com/alecthomas/gometalinter/releases/download/v2.0.11/gometalinter-2.0.11-linux-amd64.tar.gz
        tar -xvf gometalinter-2.0.11-linux-amd64.tar.gz
        break
    fi
done
cd $CURRENT_DIR