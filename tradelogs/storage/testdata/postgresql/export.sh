#!/bin/bash
# -*- firestarter: "shfmt -i 4 -ci -w %p" -*-
# Usage: ./export.sh <core-url> <core-singing-key>

set -euo pipefail

readonly from_block=6494037
readonly to_block=6494137
readonly postgresql_image='postgres:10-alpine'
readonly container_name='trade-logs-sample-data-exporter'
readonly data_dir='/tmp/trade-logs-sample-data-exporter'


docker run --rm --detach --name "$container_name" \
    --volume "$data_dir:/var/lib/postgresql/data" \
    --publish '127.0.0.1:5432:5432' \
    -e "POSTGRES_PASSWORD=reserve_stats" -e "POSTGRES_USER=reserve_stats" -e "POSTGRES_DB=reserve_stats" \
    "$postgresql_image"

sleep 60

pushd ../../../cmd/trade-logs-crawler
echo "start build"
go build
echo "start"
./trade-logs-crawler --db-engine=postgres --postgres-host=127.0.0.1 --postgres-port=5432 \
    --from-block "$from_block" --to-block "$to_block" $TL_CRAWLER_PARAMS
echo "done"
popd

docker exec "$container_name" bin/bash -c "pg_dump -U reserve_stats --attribute-inserts --data-only reserve_stats > /var/lib/postgresql/data/export.sql"

cp "$data_dir/export.sql" export.sql

docker kill "$container_name"
sudo rm -rf "$data_dir"
