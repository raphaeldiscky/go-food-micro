#!/bin/sh

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
# Change to project root directory (one level up from scripts)
cd "$SCRIPT_DIR/.." || exit 1

# Check if required parameters are provided
if [ -z "$1" ]; then
    echo "❌ Error: Test file name is required"
    echo "Usage: ./scripts/test-file.sh <test_file_name> [test_type]"
    echo "Example: ./scripts/test-file.sh product_deleted_test.go integration"
    exit 1
fi

TEST_FILE="$1"
TEST_TYPE="${2:-unit}"  # Default to unit if not specified

# Validate test type
if [ "$TEST_TYPE" != "unit" ] && [ "$TEST_TYPE" != "integration" ] && [ "$TEST_TYPE" != "e2e" ]; then
    echo "❌ Error: Invalid test type. Must be one of: unit, integration, e2e"
    exit 1
fi

# Function to find test file recursively
find_test_file() {
    service="$1"
    test_type="$2"
    test_file="$3"
    search_path=""

    # Handle pkg service differently since it's in internal/pkg not internal/services/pkg
    if [ "$service" = "pkg" ]; then
        search_path="internal/pkg"
    else
        case "$test_type" in
            "integration")
                search_path="internal/services/$service/test/integration"
                ;;
            "e2e")
                search_path="internal/services/$service/test/e2e"
                ;;
            *)
                search_path="internal/services/$service/internal"
                ;;
        esac
    fi

    # Use find to search recursively
    find "$search_path" -name "$test_file" -type f 2>/dev/null | head -n 1
}

# Search for the test file in each service
FOUND_FILE=""
SERVICE=""

for service in catalogwriteservice catalogreadservice orderservice pkg; do
    FOUND_FILE=$(find_test_file "$service" "$TEST_TYPE" "$TEST_FILE")
    if [ -n "$FOUND_FILE" ]; then
        SERVICE="$service"
        break
    fi
done

if [ -z "$FOUND_FILE" ]; then
    echo "❌ Error: Could not find test file '$TEST_FILE' in any service directory"
    echo "Make sure the file exists and test type is correct"
    exit 1
fi

echo "🧪 Running $TEST_TYPE test for file: $FOUND_FILE"
echo "📦 Service: $SERVICE"

# Get the directory of the test file
TEST_DIR=$(dirname "$FOUND_FILE")

# Run the test with appropriate tags from the test file's directory
# We run the entire package since Go test files need access to the whole package
case "$TEST_TYPE" in
    "integration")
        cd "$TEST_DIR" && go test -v -tags=integration .
        ;;
    "e2e")
        cd "$TEST_DIR" && go test -v -tags=e2e .
        ;;
    *)
        cd "$TEST_DIR" && go test -v -tags=unit .
        ;;
esac
