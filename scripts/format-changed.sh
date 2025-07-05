#!/bin/bash

# Script to format only changed Go files in the git staging area
set -e

echo "🔍 Identifying changed Go files..."

# Get all staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -z "$STAGED_GO_FILES" ]; then
    echo "✅ No Go files to format"
    exit 0
fi

echo "📝 Found staged Go files:"
echo "$STAGED_GO_FILES"

# Group files by service
PKG_FILES=""
CATALOG_WRITE_FILES=""
CATALOG_READ_FILES=""
ORDER_FILES=""

for file in $STAGED_GO_FILES; do
    if [[ $file == internal/pkg/* ]]; then
        PKG_FILES="$PKG_FILES $file"
    elif [[ $file == internal/services/catalogwriteservice/* ]]; then
        CATALOG_WRITE_FILES="$CATALOG_WRITE_FILES $file"
    elif [[ $file == internal/services/catalogreadservice/* ]]; then
        CATALOG_READ_FILES="$CATALOG_READ_FILES $file"
    elif [[ $file == internal/services/orderservice/* ]]; then
        ORDER_FILES="$ORDER_FILES $file"
    fi
done

# Function to format files for a specific service
format_service_files() {
    local service_name="$1"
    local files="$2"
    local service_path="$3"
    
    if [ -n "$files" ]; then
        echo "🧹 Formatting $service_name files..."
        cd "$service_path"
        
        # Run formatting tools on the specific files
        for file in $files; do
            # Get relative path from service directory
            rel_file=$(echo "$file" | sed "s|$service_path/||")
            echo "  📄 Formatting $rel_file"
            
            # Run golangci-lint fmt on the specific file
            golangci-lint fmt "$rel_file" || true
            
            # Run goimports to organize imports
            goimports -w "$rel_file" || true
            
            # Run gofmt to format code
            gofmt -w "$rel_file" || true
        done
        
        cd - > /dev/null
        echo "✅ $service_name formatting completed"
    fi
}

# Format files for each service if there are any
if [ -n "$PKG_FILES" ]; then
    format_service_files "pkg" "$PKG_FILES" "./internal/pkg"
fi

if [ -n "$CATALOG_WRITE_FILES" ]; then
    format_service_files "catalogwriteservice" "$CATALOG_WRITE_FILES" "./internal/services/catalogwriteservice"
fi

if [ -n "$CATALOG_READ_FILES" ]; then
    format_service_files "catalogreadservice" "$CATALOG_READ_FILES" "./internal/services/catalogreadservice"
fi

if [ -n "$ORDER_FILES" ]; then
    format_service_files "orderservice" "$ORDER_FILES" "./internal/services/orderservice"
fi

echo "✅ Formatting of changed files completed!" 