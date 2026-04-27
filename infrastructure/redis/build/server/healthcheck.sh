#!/bin/sh
set -e

if [ "$(redis-cli --no-auth-warning -u redis://default:healthcheck@127.0.0.1:6379 PING)" != "PONG" ]; then
  exit 1
fi
