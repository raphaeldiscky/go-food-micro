set dotenv

# Install development tools
install-tools:
    #!/usr/bin/env bash
    ./scripts/install-tools.sh

# Run services
run-catalogs-write-service:
    #!/usr/bin/env bash
    ./scripts/run.sh catalogwriteservice

run-catalog-read-service:
    #!/usr/bin/env bash
    ./scripts/run.sh catalogreadservice

run-order-service:
    #!/usr/bin/env bash
    ./scripts/run.sh orderservice

# Build services
build *services="pkg catalogwriteservice catalogreadservice orderservice":
    #!/usr/bin/env bash
    for service in {{services}}; do
        ./scripts/build.sh $service
    done

# Update dependencies
update-dependencies *services="pkg catalogwriteservice catalogreadservice orderservice":
    #!/usr/bin/env bash
    for service in {{services}}; do
        ./scripts/update-dependencies.sh $service
    done

# Install dependencies
install-dependencies *services="pkg catalogwriteservice catalogreadservice orderservice":
    #!/usr/bin/env bash
    for service in {{services}}; do
        ./scripts/install-dependencies.sh $service
    done

# Docker compose commands
docker-compose-infra-up:
    #!/usr/bin/env bash
    docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml up --build -d

docker-compose-infra-down:
    #!/usr/bin/env bash
    docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml down

# Generate OpenAPI specs
openapi *services="catalogwriteservice catalogreadservice orderservice":
    #!/usr/bin/env bash
    for service in {{services}}; do
        ./scripts/openapi.sh $service
    done

# Generate protobuf files
proto *services="catalogwriteservice orderservice":
    #!/usr/bin/env bash
    for service in {{services}}; do
        ./scripts/proto.sh $service
    done

# Run tests
unit-test *services="catalogreadservice catalogwriteservice orderservice":
    #!/usr/bin/env bash
    for service in {{services}}; do
        ./scripts/test.sh $service unit
    done

integration-test *services="catalogreadservice catalogwriteservice orderservice":
    #!/usr/bin/env bash
    for service in {{services}}; do
        ./scripts/test.sh $service integration
    done

e2e-test *services="catalogreadservice catalogwriteservice orderservice":
    #!/usr/bin/env bash
    for service in {{services}}; do
        ./scripts/test.sh $service e2e
    done

# Code quality
format *services="catalogwriteservice catalogreadservice orderservice pkg":
    #!/usr/bin/env bash
    for service in {{services}}; do
        ./scripts/format.sh $service
    done

lint *services="catalogwriteservice catalogreadservice orderservice pkg":
    #!/usr/bin/env bash
    for service in {{services}}; do
        ./scripts/lint.sh $service
    done

# Database migrations
go-migrate:
    #!/usr/bin/env bash
    ./scripts/go-migrate.sh -p ./internal/services/catalogwriteservice/db/migrations/go-migrate -c create -n create_product_table
    ./scripts/go-migrate.sh -p ./internal/services/catalogwriteservice/db/migrations/go-migrate -c up -o postgres://postgres:postgres@localhost:5432/catalogs_write_service?sslmode=disable
    ./scripts/go-migrate.sh -p ./internal/services/catalogwriteservice/db/migrations/go-migrate -c down -o postgres://postgres:postgres@localhost:5432/catalogs_write_service?sslmode=disable

goose-migrate:
    #!/usr/bin/env bash
    ./scripts/goose-migrate.sh -p ./internal/services/catalogwriteservice/db/migrations/goose-migrate -c create -n create_product_table
    ./scripts/goose-migrate.sh -p ./internal/services/catalogwriteservice/db/migrations/goose-migrate -c up -o "user=postgres password=postgres dbname=catalogs_write_service sslmode=disable"
    ./scripts/goose-migrate.sh -p ./internal/services/catalogwriteservice/db/migrations/goose-migrate -c down -o "user=postgres password=postgres dbname=catalogs_write_service sslmode=disable"

atlas:
    #!/usr/bin/env bash
    ./scripts/atlas-migrate.sh -c gorm-sync -p "./internal/services/catalogwriteservice"
    ./scripts/atlas-migrate.sh -c apply -p "./internal/services/catalogwriteservice" -o "postgres://postgres:postgres@localhost:5432/catalogs_write_service?sslmode=disable"

# Code analysis
cycle-check:
    #!/usr/bin/env bash
    cd internal/pkg && goimportcycle -dot imports.dot dot -Tpng -o cycle/pkg.png imports.dot
    cd internal/services/catalogwriteservice && goimportcycle -dot imports.dot dot -Tpng -o cycle/catalogwriteservice.png imports.dot
    cd internal/services/catalogreadservice && goimportcycle -dot imports.dot dot -Tpng -o cycle/catalogreadservice.png imports.dot
    cd internal/services/orderservice && goimportcycle -dot imports.dot dot -Tpng -o cycle/orderservice.png imports.dot

# Generate mocks
pkg-mocks:
    #!/usr/bin/env bash
    cd internal/pkg/es && mockery --output mocks --all
    cd internal/pkg/core/serializer && mockery --output mocks --all
    cd internal/pkg/core/messaging && mockery --output mocks --all

services-mocks:
    #!/usr/bin/env bash
    cd internal/services/catalogwriteservice && mockery --output mocks --all
    cd internal/services/catalogreadservice && mockery --output mocks --all
    cd internal/services/orderservice && mockery --output mocks --all
