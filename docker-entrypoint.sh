#!/bin/sh
set -eu

./app serve --http=0.0.0.0:8090 &
backend_pid=$!

node server.js &
frontend_pid=$!

shutdown() {
  kill -TERM "$backend_pid" "$frontend_pid" 2>/dev/null || true
  wait "$backend_pid" "$frontend_pid" 2>/dev/null || true
}

trap shutdown INT TERM
wait "$backend_pid" "$frontend_pid"
