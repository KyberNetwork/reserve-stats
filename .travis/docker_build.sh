#!/bin/bash

set -euo pipefail

readonly module=${MODULE:-}

declare -a services_list=()

if [[ "$module" == "reserverates" ]]; then
    services_list=("reserve-rates-api" "reserve-rates-crawler")
elif [[ "$module" == "tradelogs" ]]; then
    services_list=("trade-logs-api" "trade-logs-crawler")
elif [[ "$module" == "users" ]]; then
    services_list=("users-api")
else
    exit 0
fi

for service in ${services_list[@]}; do
    docker build -f docker-files/Dockerfile.$service -t kybernetwork/kyber-stats-$service:$TRAVIS_COMMIT .
done
