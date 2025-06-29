#!/bin/bash

# Update all ginkgo v1 imports to v2 in Go files
set -e

echo "Updating ginkgo imports from v1 to v2..."

# Update dot imports
find internal/services -name "*.go" -type f -exec sed -i 's|"github\.com/onsi/ginkgo"|"github.com/onsi/ginkgo/v2"|g' {} \;

# Update named imports
find internal/services -name "*.go" -type f -exec sed -i 's|ginkgo "github\.com/onsi/ginkgo"|ginkgo "github.com/onsi/ginkgo/v2"|g' {} \;

echo "Ginkgo imports updated successfully!"
echo "Now running go mod tidy for all services..."

# Run go mod tidy for all services to clean up dependencies
cd internal/services/orderservice && go mod tidy
cd ../catalogwriteservice && go mod tidy  
cd ../catalogreadservice && go mod tidy

echo "All done! Ginkgo v2 migration completed."