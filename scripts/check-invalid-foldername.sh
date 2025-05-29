#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Function to check if a string is lowercase and single word (allowing numbers)
is_valid_folder_name() {
    local name=$1
    # Check if name contains only lowercase letters and numbers, no spaces or special characters
    if [[ "$name" =~ ^[a-z0-9]+$ ]]; then
        return 0
    else
        return 1
    fi
}

# Directories to exclude from checking
EXCLUDE_DIRS=(
    "."
    ".git"
    "vendor"
    "node_modules"
    ".idea"
    ".vscode"
    "bin"
    "dist"
    "build"
    "third_party"
    "docker-compose"
    "go-migrate"
    "goose-migrate"
)

# Function to check if directory should be excluded
should_exclude() {
    local dir=$1
    # Get just the directory name without the path
    local dirname=$(basename "$dir")
    
    # Check for exact matches first
    for exclude in "${EXCLUDE_DIRS[@]}"; do
        if [[ "$dirname" == "$exclude" ]]; then
            return 0
        fi
    done
    
    # Check for path components
    for exclude in "${EXCLUDE_DIRS[@]}"; do
        if [[ "$dir" == *"/$exclude"* ]]; then
            return 0
        fi
    done
    return 1
}

# Find all directories and check their names
echo "Checking for invalid folder names..."
echo "Invalid folder names should be single word, lowercase letters and numbers only"
echo "----------------------------------------"

found_invalid=0

# Find all directories and process them
while IFS= read -r -d '' dir; do
    # Get just the directory name without the path
    dirname=$(basename "$dir")
    
    # Skip if directory should be excluded
    if should_exclude "$dir"; then
        continue
    fi
    
    # Check if the directory name is valid
    if ! is_valid_folder_name "$dirname"; then
        echo -e "${RED}Invalid folder name found:${NC} $dir"
        echo "  - Should be a single word in lowercase letters and numbers only"
        found_invalid=1
    fi
done < <(find . -type d -not -path "*/\.*" -print0)

if [ $found_invalid -eq 0 ]; then
    echo -e "${GREEN}All folder names are valid!${NC}"
    exit 0
else
    echo -e "\n${RED}Found invalid folder names. Please fix them to follow Go conventions.${NC}"
    exit 1
fi
