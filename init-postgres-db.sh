#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE users;
    CREATE DATABASE "price_analytics";
    CREATE DATABASE "accounting";
    CREATE DATABASE "app-names";
    CREATE DATABASE "cex_trades";
    CREATE DATABASE "cex_withdrawals";
    CREATE DATABASE "listed_tokens";
    CREATE DATABASE "reserve_rates";
    CREATE DATABASE "transactions";
    CREATE DATABASE "reserve_addresses";
EOSQL