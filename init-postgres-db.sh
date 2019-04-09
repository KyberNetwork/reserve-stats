#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE users;
    CREATE DATABASE "price-analytics";
    CREATE DATABASE "accounting";
    CREATE DATABASE "app-names";
EOSQL