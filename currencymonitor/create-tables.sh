#!/bin/bash

set -e
set -x

docker-compose exec scylladb cqlsh -e "DROP KEYSPACE IF EXISTS currencymonitor;"
docker-compose exec scylladb cqlsh -e "CREATE KEYSPACE currencymonitor WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};"
docker-compose exec scylladb cqlsh -e "CREATE TABLE currencymonitor.monitoring (id TEXT, status tinyint, period TEXT, frequency TEXT, data TEXT, PRIMARY KEY (id));"
docker-compose exec scylladb cqlsh -e "CREATE INDEX IF NOT EXISTS monitoring_status ON currencymonitor.monitoring (status);"