#!/bin/bash
# -*- firestarter: "shfmt -i 4 -ci -w %p" -*-
# Usage: ./export.sh <core-url> <core-singing-key>

set -euo pipefail

readonly core_url=${1:-}
readonly core_signing_key=${2:-}

readonly from_block=6494037
readonly to_block=6494137
readonly influxdb_image='influxdb:1.7.1-alpine'
readonly container_name='trade-logs-sample-data-exporter'
readonly data_dir='/tmp/trade-logs-sample-data-exporter'


docker run --rm --detach --name "$container_name" \
    --volume "$data_dir:/var/lib/influxdb/" \
    --publish '127.0.0.1:8087:8086' \
    "$influxdb_image"

sleep 5

pushd ../../cmd/trade-logs-crawler
go build
./trade-logs-crawler --influxdb-endpoint=http://127.0.0.1:8087 \
    --from-block "$from_block" --to-block "$to_block"
popd

docker exec "$container_name" influx_inspect export -datadir /var/lib/influxdb/data \
    -waldir /var/lib/influxdb/wal -out /var/lib/influxdb/data/export.dat \
    -database trade_logs -retention autogen

docker exec "$container_name" sed -i -e '/^#\|^CREATE DATABASE/d' /var/lib/influxdb/data/export.dat 

cp "$data_dir/data/export.dat" export.dat

docker kill "$container_name"
sudo rm -rf "$data_dir"
