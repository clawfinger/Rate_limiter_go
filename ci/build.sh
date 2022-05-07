#!/bin/bash

cd /service &&
go build -v -o ./bin ./cmd/limiter &&
go build -v -o ./bin ./cmd/cli