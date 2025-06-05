#!/bin/bash

# https://blog.devgenius.io/go-golang-testing-tools-tips-to-step-up-your-game-4ed165a5b3b5

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

readonly service="$1"
readonly type="$2"

# Function to run tests for a specific service
run_service_tests() {
    local service_path="$1"
    local service_name=$(basename "$service_path")
    local test_type="$2"
    
    echo "[$service_name] Starting tests..."
    cd "$service_path"
    
    # Function to run tests and capture their exit status
    run_tests() {
        local test_type="$1"
        echo "[$service_name] Starting $test_type tests..."
        go test -tags="$test_type" -timeout=30m -v -count=1 -p=1 -parallel=1 ./...
        local status=$?
        if [ $status -ne 0 ]; then
            echo "[$service_name] $test_type tests failed with exit status $status"
            exit $status
        fi
        echo "[$service_name] $test_type tests completed successfully"
    }

    if [ "$test_type" = "load-test" ]; then
        echo "[$service_name] Running load tests..."
        k6 run ./load_tests/script.js --insecure-skip-tls-verify
    elif [ -z "$test_type" ]; then
        # Run all test types concurrently
        echo "[$service_name] Running all test types concurrently..."
        
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
        
        echo "[$service_name] All tests completed successfully!"
    else
        # Run specific test type
        run_tests "$test_type"
    fi
    
    cd - > /dev/null  # Return to previous directory
}

# If a specific service is provided, run only for that service
if [ -n "$service" ]; then
    if [ "$service" = "pkg" ]; then
        run_service_tests "./internal/pkg" "$type"
    else
        run_service_tests "./internal/services/$service" "$type"
    fi
else
    # Run for all services concurrently
    echo "Running tests for all services concurrently..."
    
    # Store PIDs in a space-separated string instead of an array
    pids=""
    
    # Run for pkg
    run_service_tests "./internal/pkg" "$type" &
    pids="$pids $!"
    
    # Run for each service in services directory
    for service_dir in ./internal/services/*/; do
        if [ -d "$service_dir" ]; then
            run_service_tests "$service_dir" "$type" &
            pids="$pids $!"
        fi
    done
    
    # Wait for all background processes to complete
    for pid in $pids; do
        wait $pid
    done
    
    echo "All services tests completed!"
fi


