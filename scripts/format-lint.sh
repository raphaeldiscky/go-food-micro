#!/bin/bash

# ref: https://blog.devgenius.io/sort-go-imports-acb76224dfa7
# https://yolken.net/blog/cleaner-go-code-golines

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

readonly service="$1"

# Function to run format-lint for a specific service
run_format_lint() {
    local service_path="$1"
    local service_name=$(basename "$service_path")
    echo "Starting format-lint for $service_name..."
    
    cd "$service_path"
    
    # Function to run a command and capture its exit status
    run_cmd() {
        local cmd="$1"
        local name="$2"
        echo "[$service_name] Starting $name..."
        $cmd
        local status=$?
        if [ $status -ne 0 ]; then
            echo "[$service_name] $name failed with exit status $status"
            exit $status
        fi
        echo "[$service_name] $name completed successfully"
    }

    # Run linting tools concurrently
    run_cmd "revive -config revive-config.toml -formatter friendly ./..." "revive" &
    revive_pid=$!

    run_cmd "staticcheck ./..." "staticcheck" &
    staticcheck_pid=$!

    run_cmd "golangci-lint run --fix ./..." "golangci-lint run" &
    golangci_run_pid=$!

    run_cmd "golangci-lint fmt ./..." "golangci-lint fmt" &
    golangci_fmt_pid=$!

    # Wait for all background processes to complete
    wait $revive_pid $staticcheck_pid $golangci_run_pid $golangci_fmt_pid

    # Run errcheck if available
    if command -v errcheck >/dev/null 2>&1; then
        echo "[$service_name] Running errcheck..."
        errcheck ./...
    fi

    echo "[$service_name] All linting and formatting completed!"
    cd - > /dev/null  # Return to previous directory
}

# If a specific service is provided, run only for that service
if [ -n "$service" ]; then
    if [ "$service" = "pkg" ]; then
        run_format_lint "./internal/pkg"
    else
        run_format_lint "./internal/services/$service"
    fi
else
    # Run for all services concurrently
    echo "Running format-lint for all services concurrently..."
    
    # Store PIDs in a space-separated string instead of an array
    pids=""
    
    # Run for pkg
    run_format_lint "./internal/pkg" &
    pids="$pids $!"
    
    # Run for each service in services directory
    for service_dir in ./internal/services/*/; do
        if [ -d "$service_dir" ]; then
            run_format_lint "$service_dir" &
            pids="$pids $!"
        fi
    done
    
    # Wait for all background processes to complete
    for pid in $pids; do
        wait $pid
    done
    
    echo "All services format-lint completed!"
fi
