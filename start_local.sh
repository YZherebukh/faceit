#!/bin/bash
docker-compose -f docker/postgres/docker-compose.yaml up -d

# apt-get install build-essential -y
# apt-get install curl -y
curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
chmod +x /usr/local/bin/dbmate

make migrate-up

export SERVICE_PORT_ENV=":8080"
export VERSION_API_ENV="v1"
export DB_NAME_ENV="test_database"
export DB_PORT_ENV="5432"
export DB_HOST_ENV="localhost"
export DB_USER_NAME_ENV="postgres"
export DB_PASSWORD_ENV="abcd"
export DB_SSLMODE_ENV="disable"
export DB_URL_PATTERN_ENV="host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"
export LOGGER_PATH_ENV="./code_test.log"
export LOGGER_VERBOSE_ENV="true"
export LOGGER_SYSTEM_LOG_ENV="false"
export NOTIFIER_CONSUMERS_CREATE_ENV="" 
export NOTIFIER_CONSUMERS_UPDATE_ENV=""
export NOTIFIER_CONSUMERS_DELETE_ENV=""
export NOTIFIER_TIMEOUT_ENV=1
export NOTIFIER_CLIENT_MAX_RETRY_ENV=3
export NOTIFIER_TIMEOUT_INCREACE_ENV=3

go build && ./test 
