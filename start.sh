#!/bin/bash

export DATABASE_URL=postgres://postgres:abcd@mypostgres:5432/test_database?sslmode=disable

curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
chmod +x /usr/local/bin/dbmate

echo $DATABASE_URL

make migrate-up

./test