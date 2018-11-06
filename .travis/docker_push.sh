#!/bin/bash
# -*- firestarter: "shfmt -i 4 -ci -w %p" -*-

set -euo pipefail

readonly docker_password=${DOCKER_PASSWORD:-}
readonly build_part=${BUILD_PART:-}

push() {
    for service in "${@}"; do
        local docker_repository="kybernetwork/kyber-stats-$service"
        docker tag "$docker_repository:$TRAVIS_COMMIT" "$docker_repository:$TRAVIS_BRANCH"
        if [[ -n "$TRAVIS_TAG" ]]; then
            docker tag "$docker_repository:$TRAVIS_COMMIT" "$docker_repository:$TRAVIS_TAG"
        fi
        docker push "$docker_repository"
    done
}

echo "$docker_password" | docker login -u "$DOCKER_USERNAME" --password-stdin

case "$build_part" in
    1)
        push reserverates reserve-rates-api reserve-rates-crawler users-api gateway
        ;;
    2)
        push trade-logs-api trade-logs-crawler price-analytics-api
        ;;
esac
