#!/bin/bash

# ref: https://blog.devgenius.io/sort-go-imports-acb76224dfa7
# https://yolken.net/blog/cleaner-go-code-golines

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

readonly service="$1"

if [ "$service" = "pkg" ]; then
      cd "./internal/pkg"
# Check if input is not empty or null
elif [ -n "$service"  ]; then
    cd "./internal/services/$service"
fi

# Function to run a command and capture its exit status
run_cmd() {
    local cmd="$1"
    local name="$2"
    echo "Starting $name..."
    $cmd
    local status=$?
    if [ $status -ne 0 ]; then
        echo "$name failed with exit status $status"
        exit $status
    fi
    echo "$name completed successfully"
}

# Run linting tools concurrently
echo "Starting parallel linting..."

# Run revive in background
run_cmd "revive -config revive-config.toml -formatter friendly ./..." "revive" &
revive_pid=$!

# Run staticcheck in background
run_cmd "staticcheck ./..." "staticcheck" &
staticcheck_pid=$!

# Run golangci-lint commands in background
run_cmd "golangci-lint run --fix ./..." "golangci-lint run" &
golangci_run_pid=$!

run_cmd "golangci-lint fmt ./..." "golangci-lint fmt" &
golangci_fmt_pid=$!

# Wait for all background processes to complete
wait $revive_pid $staticcheck_pid $golangci_run_pid $golangci_fmt_pid

# Run errcheck if available (keeping this sequential since it's optional)
if command -v errcheck >/dev/null 2>&1; then
    echo "Running errcheck..."
    errcheck ./...
fi

echo "All linting and formatting completed!"
