#!/bin/bash

set -euo pipefail

readonly docker_password=${DOCKER_PASSWORD:-}
readonly module=${MODULE:-}

declare -a docker_repositories_list=()

if [[ "$module" == "reserverates" ]]; then
    docker_repositories_list=("kyber-stats-reserve-rates-api" "kyber-stats-reserve-rates-crawler")
elif [[ "$module" == "tradelogs" ]]; then
    docker_repositories_list=("kyber-stats-trade-logs-api" "kyber-stats-trade-logs-crawler")
elif [[ "$module" == "users" ]]; then
    docker_repositories_list=("kyber-stats-users-api")
elif [[ "$module" == "others" ]]; then
    exit 0
else
    echo "Module $module is not an valid module"
    exit 1
fi

if [[ -z "$docker_password" ]]; then
    echo 'DOCKER_PASSWORD is not available, aborting.'
    exit 1
fi

echo "$docker_password" | docker login -u "$DOCKER_USERNAME" --password-stdin

for docker_repository in ${docker_repositories_list[@]}; do
    docker tag "kybernetwork/$docker_repository:$TRAVIS_COMMIT" "kybernetwork/$docker_repository:$TRAVIS_BRANCH"
    if [[ -n "$TRAVIS_TAG" ]]; then
        docker tag "kybernetwork/$docker_repository:$TRAVIS_COMMIT" "kybernetwork/$docker_repository:$TRAVIS_TAG"
    fi
    docker push "kybernetwork/$docker_repository"
done
