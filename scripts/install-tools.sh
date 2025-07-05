#!/bin/bash

# https://github.com/actions/runner-images/blob/main/images/linux/Ubuntu2204-Readme.md

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

# `go install package@version` command works directly when we specified exact version, elsewhere it needs a `go.mod` and specifying corresponding version for each package

go install github.com/samlitowitz/goimportcycle/cmd/goimportcycle@v1.0.9

# https://github.com/incu6us/goimports-reviser
go install -v github.com/incu6us/goimports-reviser/v3@v3.6.5

# https://github.com/daixiang0/gci
go install github.com/daixiang0/gci@v0.13.4

# https://pkg.go.dev/golang.org/x/tools/cmd/goimports
go install golang.org/x/tools/cmd/goimports@v0.34.0

# https://github.com/mvdan/gofumpt
go install mvdan.cc/gofumpt@v0.8.0

# https://github.com/segmentio/golines
go install github.com/segmentio/golines@v0.12.2

# https://golangci-lint.run/welcome/install/
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.2.0

# https://github.com/mgechev/revive
go install github.com/mgechev/revive@v1.10.0

# https://github.com/dominikh/go-tools
go install honnef.co/go/tools/cmd/staticcheck@2025.1.1

# https://dev.to/techschoolguru/how-to-define-a-protobuf-message-and-generate-go-code-4g4e
# https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

# migration tools
# https://github.com/pressly/goose
go install github.com/pressly/goose/v3/cmd/goose@v3.24.3
# https://github.com/golang-migrate/migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

# https://github.com/swaggo/swag/
# https://github.com/swaggo/swag/issues/817
# Using v1.16.4 for stability with generic types
go install github.com/swaggo/swag/cmd/swag@v1.16.4

# https://github.com/deepmap/oapi-codegen
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

# https://github.com/OpenAPITools/openapi-generator-cli
npm install -g @openapitools/openapi-generator-cli@2.14.0

# https://vektra.github.io/mockery/latest/installation/
go install github.com/vektra/mockery/v2@v2.46.0

# https://github.com/onsi/ginkgo
go install github.com/onsi/ginkgo/v2/ginkgo@v2.23.4

# https://github.com/bufbuild/buf
go install github.com/bufbuild/buf/cmd/buf@v1.55.1

OS="$(uname -s)"

if [[ "$OS" == "Linux" ]]; then
    # https://github.com/ariga/atlas#quick-installation
    curl -sSf https://atlasgo.sh | sh

    # https://k6.io/docs/get-started/installation/
    sudo gpg -k
    sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
    echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
    sudo apt-get update
    sudo apt-get install k6

    # https://grpc.io/docs/protoc-installation/
    sudo apt install -y protobuf-compiler
else
    echo "Unsupported operating system: $OS"
    exit 1
fi