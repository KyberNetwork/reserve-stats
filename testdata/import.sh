#!/bin/bash

set -euo pipefail

# /usr/bin/pg_dump --dbname=accounting --clean --if-exists --file=accounting.sql --username=reserve_stats --host=localhost --port=5432
echo "importing accounting test database"
export PGPASSWORD=reserve_stats
psql -h 127.0.0.1 -U reserve_stats -d accounting < ./accounting.sql