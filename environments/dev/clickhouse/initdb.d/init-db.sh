#!/bin/bash

set -e

database=${GLE_CLICKHOUSE_DATABASE:-'gitlab_ci'}
username=${GLE_CLICKHOUSE_USER:-'gitlab_exporter'}
password=${GLE_CLICKHOUSE_PASSWORD:-'gitlab_exporter'}

clickhouse client -n <<-EOSQL
    CREATE DATABASE IF NOT EXISTS ${database};

    CREATE USER IF NOT EXISTS ${username}
        IDENTIFIED WITH sha256_password BY '${password}'
        HOST ANY
        ;

    GRANT ALL ON ${database}.* TO ${username};
EOSQL
