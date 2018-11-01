#!/bin/bash

set -euo pipefail

readonly module=${MODULE:-}

declare -a services_list=()

gometalinter --config=gometalinter.json ./...
go build -v -mod=vendor ./...
go test -v -mod=vendor ./...

if [[ "$module" =~ (reserverates|tradelogs|users) ]]; then
    echo "Testing $module module"
    (cd $module; go build -v ./...; go test -v ./...)
elif [[ "$module" == "others" ]]; then
    go build -v $(go list ./... | grep -v "github.com/KyberNetwork/reserve-stats/\(reserverates\|tradelogs\|users\)")
    exit 0
else
    echo "Module $module is not an valid module"
    exit 1
fi

if [[ "$module" == "reserverates" ]]; then
    services_list=("reserve-rates-api" "reserve-rates-crawler")
elif [[ "$module" == "tradelogs" ]]; then
    services_list=("trade-logs-api" "trade-logs-crawler")
elif [[ "$module" == "users" ]]; then
    services_list=("users-api")
elif [[ "$module" == "others" ]]; then
    exit 0
else
    echo "Module $module is not an valid module"
    exit 1
fi

for service in ${services_list[@]}; do
    docker build -f docker-files/Dockerfile.$service -t kybernetwork/kyber-stats-$service:$TRAVIS_COMMIT .
done
