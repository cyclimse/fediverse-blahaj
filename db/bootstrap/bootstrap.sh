#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER fediverse WITH PASSWORD 'fediverse';
	CREATE DATABASE fediverse;
	GRANT ALL PRIVILEGES ON DATABASE fediverse TO fediverse;

	CREATE DATABASE fediversedev;
	GRANT ALL PRIVILEGES ON DATABASE fediversedev TO fediverse;
EOSQL

# Install extensions

for db in fediverse fediversedev; do
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$db" <<-EOSQL
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
EOSQL
done
