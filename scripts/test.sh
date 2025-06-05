#!/bin/bash

# https://blog.devgenius.io/go-golang-testing-tools-tips-to-step-up-your-game-4ed165a5b3b5

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

readonly service="$1"
readonly type="$2"

if [ "$service" = "pkg" ]; then
      cd "./internal/pkg/$service"
else
    cd "./internal/services/$service"
fi

# Function to run tests and capture their exit status
run_tests() {
    local test_type="$1"
    echo "Starting $test_type tests..."
    go test -tags="$test_type" -timeout=30m -v -count=1 -p=1 -parallel=1 ./...
    local status=$?
    if [ $status -ne 0 ]; then
        echo "$test_type tests failed with exit status $status"
        exit $status
    fi
    echo "$test_type tests completed successfully"
}

if [ "$type" = "load-test" ]; then
    # go run ./cmd/app/main.go
    k6 run ./load_tests/script.js --insecure-skip-tls-verify
elif [ -z "$type" ]; then
    # Run all test types concurrently when no specific type is provided
    echo "Running all test types concurrently..."
    
    # Run unit tests in background
    run_tests "unit" &
    unit_pid=$!
    
    # Run integration tests in background
    run_tests "integration" &
    integration_pid=$!
    
    # Run e2e tests in background
    run_tests "e2e" &
    e2e_pid=$!
    
    # Wait for all test processes to complete
    wait $unit_pid $integration_pid $e2e_pid
    
    echo "All tests completed successfully!"
else
    # Run specific test type
    run_tests "$type"
fi


