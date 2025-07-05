#!/bin/bash

# Safe Dependency Update Script for Go Microservices
# This script provides a comprehensive approach to updating Go dependencies safely

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BACKUP_DIR="./backup-$(date +%Y%m%d-%H%M%S)"
SERVICES=("pkg" "catalogwriteservice" "catalogreadservice" "orderservice")

# Logging function
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        error "Go is not installed or not in PATH"
        exit 1
    fi
    log "Go version: $(go version)"
}

# Create backup of go.mod and go.sum files
create_backup() {
    log "Creating backup in $BACKUP_DIR"
    mkdir -p "$BACKUP_DIR"
    
    for service in "${SERVICES[@]}"; do
        if [ "$service" = "pkg" ]; then
            service_path="./internal/pkg"
        else
            service_path="./internal/services/$service"
        fi
        
        if [ -d "$service_path" ]; then
            mkdir -p "$BACKUP_DIR/$service"
            if [ -f "$service_path/go.mod" ]; then
                cp "$service_path/go.mod" "$BACKUP_DIR/$service/"
            fi
            if [ -f "$service_path/go.sum" ]; then
                cp "$service_path/go.sum" "$BACKUP_DIR/$service/"
            fi
        fi
    done
    success "Backup created successfully"
}

# Check for outdated dependencies
check_outdated() {
    log "Checking for outdated dependencies..."
    
    for service in "${SERVICES[@]}"; do
        if [ "$service" = "pkg" ]; then
            service_path="./internal/pkg"
        else
            service_path="./internal/services/$service"
        fi
        
        if [ -d "$service_path" ]; then
            log "Checking $service..."
            cd "$service_path"
            
            # Check for outdated modules
            outdated=$(go list -u -m all 2>/dev/null | grep -E "\[.*\]" || true)
            if [ -n "$outdated" ]; then
                warning "Outdated dependencies found in $service:"
                echo "$outdated"
            else
                success "$service is up to date"
            fi
            
            cd - > /dev/null
        fi
    done
}

# Update dependencies for a specific service
update_service() {
    local service=$1
    local service_path
    
    if [ "$service" = "pkg" ]; then
        service_path="./internal/pkg"
    else
        service_path="./internal/services/$service"
    fi
    
    if [ ! -d "$service_path" ]; then
        error "Service directory not found: $service_path"
        return 1
    fi
    
    log "Updating dependencies for $service..."
    cd "$service_path"
    
    # Update all dependencies
    log "Running go get -u -t -d -v ./..."
    go get -u -t -d -v ./...
    
    # Tidy up the module
    log "Running go mod tidy"
    go mod tidy
    
    # Verify the module
    log "Verifying module..."
    go mod verify
    
    cd - > /dev/null
    success "Updated dependencies for $service"
}

# Test the service after update
test_service() {
    local service=$1
    local service_path
    
    if [ "$service" = "pkg" ]; then
        service_path="./internal/pkg"
    else
        service_path="./internal/services/$service"
    fi
    
    if [ ! -d "$service_path" ]; then
        error "Service directory not found: $service_path"
        return 1
    fi
    
    log "Testing $service..."
    cd "$service_path"
    
    # Run tests
    if go test ./... -v; then
        success "Tests passed for $service"
    else
        error "Tests failed for $service"
        return 1
    fi
    
    cd - > /dev/null
}

# Build the service to check for compilation errors
build_service() {
    local service=$1
    local service_path
    
    if [ "$service" = "pkg" ]; then
        service_path="./internal/pkg"
    else
        service_path="./internal/services/$service"
    fi
    
    if [ ! -d "$service_path" ]; then
        error "Service directory not found: $service_path"
        return 1
    fi
    
    log "Building $service..."
    cd "$service_path"
    
    # Try to build
    if go build ./...; then
        success "Build successful for $service"
    else
        error "Build failed for $service"
        return 1
    fi
    
    cd - > /dev/null
}

# Main update function
update_all() {
    log "Starting safe dependency update process..."
    
    # Check prerequisites
    check_go
    
    # Create backup
    create_backup
    
    # Check current state
    check_outdated
    
    # Update each service
    for service in "${SERVICES[@]}"; do
        log "Processing $service..."
        
        if update_service "$service"; then
            if build_service "$service"; then
                if test_service "$service"; then
                    success "$service updated successfully"
                else
                    error "Tests failed for $service - consider rolling back"
                    return 1
                fi
            else
                error "Build failed for $service - consider rolling back"
                return 1
            fi
        else
            error "Update failed for $service"
            return 1
        fi
    done
    
    success "All dependencies updated successfully!"
    log "Backup available at: $BACKUP_DIR"
}

# Rollback function
rollback() {
    log "Rolling back to backup..."
    
    for service in "${SERVICES[@]}"; do
        if [ "$service" = "pkg" ]; then
            service_path="./internal/pkg"
        else
            service_path="./internal/services/$service"
        fi
        
        if [ -f "$BACKUP_DIR/$service/go.mod" ]; then
            log "Rolling back $service..."
            cp "$BACKUP_DIR/$service/go.mod" "$service_path/"
            if [ -f "$BACKUP_DIR/$service/go.sum" ]; then
                cp "$BACKUP_DIR/$service/go.sum" "$service_path/"
            fi
            cd "$service_path"
            go mod tidy
            cd - > /dev/null
            success "Rolled back $service"
        fi
    done
}

# Show usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  update     Update all dependencies (default)"
    echo "  check      Check for outdated dependencies"
    echo "  rollback   Rollback to backup"
    echo "  help       Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 update     # Update all dependencies"
    echo "  $0 check      # Check what's outdated"
    echo "  $0 rollback   # Rollback to backup"
}

# Main script logic
case "${1:-update}" in
    "update")
        update_all
        ;;
    "check")
        check_go
        check_outdated
        ;;
    "rollback")
        if [ ! -d "$BACKUP_DIR" ]; then
            error "No backup found. Please specify backup directory or run update first."
            exit 1
        fi
        rollback
        ;;
    "help"|"-h"|"--help")
        usage
        ;;
    *)
        error "Unknown option: $1"
        usage
        exit 1
        ;;
esac 