#!/bin/bash

set -euo pipefail

echo "importing accounting test database"
export PGPASSWORD=reserve_stats
psql -h 127.0.0.1 -U reserve_stats -d accounting < ./accounting.sql