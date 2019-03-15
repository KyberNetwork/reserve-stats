#!/bin/bash
# -*- firestarter: "shfmt -i 4 -ci -w %p" -*-

set -euo pipefail

readonly docker_password=${DOCKER_PASSWORD:-}
readonly build_part=${BUILD_PART:-}

push() {
    for service in "${@}"; do
        local docker_repository="kybernetwork/kyber-stats-$service"

        docker build -f "docker-files/Dockerfile.$service" -t "$docker_repository:$TRAVIS_COMMIT" .

        docker tag "$docker_repository:$TRAVIS_COMMIT" "$docker_repository:$TRAVIS_BRANCH"
        if [[ -n "$TRAVIS_TAG" ]]; then
            docker tag "$docker_repository:$TRAVIS_COMMIT" "$docker_repository:$TRAVIS_TAG"
        fi
        docker push "$docker_repository"
    done
}

push_file() {
    local config_file="$1"
    while read -r line; do
        services=($line)
        services=(${services[@]:1})
        push "${services[@]}"
    done < "$config_file"
}

echo "$docker_password" | docker login -u "$DOCKER_USERNAME" --password-stdin

case "$build_part" in
    1)
        push_file .travis/build_part_1
        ;;
    2)
        push_file .travis/build_part_2
        ;;
esac
