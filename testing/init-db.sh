#!/bin/sh
set -e

echo "Downloading Kea PostgreSQL initialization script..."
wget -q -O /tmp/dhcpdb_create.pgsql https://gitlab.isc.org/isc-projects/kea/-/raw/Kea-${KEA_VERSION}/src/share/database/scripts/pgsql/dhcpdb_create.pgsql

echo "Applying schema to database..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /tmp/dhcpdb_create.pgsql

echo "Database initialization complete"
