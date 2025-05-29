#!/bin/bash

# https://github.com/actions/runner-images/blob/main/images/linux/Ubuntu2204-Readme.md

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

# `go install package@version` command works directly when we specified exact version, elsewhere it needs a `go.mod` and specifying corresponding version for each package

go install github.com/samlitowitz/goimportcycle/cmd/goimportcycle@latest

# https://github.com/incu6us/goimports-reviser
go install -v github.com/incu6us/goimports-reviser/v3@latest

# https://github.com/daixiang0/gci
go install github.com/daixiang0/gci@latest

# https://pkg.go.dev/golang.org/x/tools/cmd/goimports
go install golang.org/x/tools/cmd/goimports@latest

# https://github.com/mvdan/gofumpt
go install mvdan.cc/gofumpt@latest

# https://github.com/segmentio/golines
go install github.com/segmentio/golines@latest

# https://golangci-lint.run/welcome/install/
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6

go install github.com/mgechev/revive@latest

# https://github.com/dominikh/go-tools
go install honnef.co/go/tools/cmd/staticcheck@latest

# https://dev.to/techschoolguru/how-to-define-a-protobuf-message-and-generate-go-code-4g4e
# https://stackoverflow.com/questions/13616033/install-protocol-buffers-on-windows
go install github.com/golang/protobuf/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# migration tools
go install github.com/pressly/goose/v3/cmd/goose@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# https://github.com/swaggo/swag/
# https://github.com/swaggo/swag/issues/817
# swag cli v1.8.3 - upper versions have some problems with generic types
go install github.com/swaggo/swag/cmd/swag@latest
# go install github.com/swaggo/swag/cmd/swag@v1.8.3

# https://github.com/deepmap/oapi-codegen
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# https://github.com/OpenAPITools/openapi-generator-cli
npm install -g @openapitools/openapi-generator-cli

# https://vektra.github.io/mockery/latest/installation/
go install github.com/vektra/mockery/v2@latest

go install github.com/onsi/ginkgo/v2/ginkgo@latest

go install github.com/bufbuild/buf/cmd/buf@latest

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
