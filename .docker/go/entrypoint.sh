#!/usr/bin/env sh
set -e

go mod download

echo
echo 'Init process done. Ready for start up.'
echo

exec "$@"
