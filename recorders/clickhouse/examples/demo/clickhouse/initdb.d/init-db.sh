#!/bin/sh

set -eu

database=${INITDB_CLICKHOUSE_DATABASE}

username=${INITDB_CLICKHOUSE_USER}
password=${INITDB_CLICKHOUSE_PASSWORD}

readonly_username=${INITDB_CLICKHOUSE_READONLY_USER}
readonly_password=${INITDB_CLICKHOUSE_READONLY_PASSWORD}

clickhouse client -n <<-EOSQL
    /*
    Create database
    */

    CREATE DATABASE IF NOT EXISTS ${database};

    /*
    Create users
    */

    CREATE USER IF NOT EXISTS ${username}
        IDENTIFIED WITH sha256_password BY '${password}'
        HOST ANY
        ;

    CREATE USER IF NOT EXISTS ${readonly_username}
        IDENTIFIED WITH sha256_password BY '${readonly_password}'
        HOST ANY
        ;
    
    /*
    Manage access
    */

    GRANT ALL ON ${database}.* TO ${username};

    GRANT SELECT,SHOW ON ${database}.* TO ${readonly_username};
EOSQL
