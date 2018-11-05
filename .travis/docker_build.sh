#!/bin/bash

set -euo pipefail

readonly module=${MODULE:-}

declare -a services_list=()

gometalinter --config=gometalinter.json ./...

if [[ "$module" =~ (reserverates|tradelogs|users|gateway) ]]; then
    echo "Testing $module module"
    (cd $module; go build -v -mod=vendor ./...; go test -v -mod=vendor ./...)
elif [[ "$module" == "others" ]]; then
    go build -v -mod=vendor $(go list ./... | grep -v "github.com/KyberNetwork/reserve-stats/\(reserverates\|tradelogs\|users\|gateway\)")
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
elif [[ "$module" == "gateway" ]]; then
    services_list=("gateway")
elif [[ "$module" == "others" ]]; then
    exit 0
else
    echo "Module $module is not an valid module"
    exit 1
fi

for service in ${services_list[@]}; do
    docker build -f docker-files/Dockerfile.$service -t kybernetwork/kyber-stats-$service:$TRAVIS_COMMIT .
done
