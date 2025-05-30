FROM golang:1.23.6-alpine AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /services/catalogs/read_service

# Install dependencies with version pinning
RUN apk --update --no-cache add ca-certificates=20240705-r0 make=4.4.1-r2 protoc=24.4-r0

# Download grpc_health_probe with newer version
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.25 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# Build Go binary
COPY Makefile go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN make init && go mod download
COPY . .
RUN make proto tidy

# Skaffold passes in debug-oriented compiler flags
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -ldflags="-w -s" -o catalogs_read .

# Deployment container
FROM gcr.io/distroless/static:nonroot

WORKDIR /read_service

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy health probe
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

# Copy binary
COPY --from=builder /services/catalogs/read_service/catalogs_read /services/catalogs-read

# Use non-root user
USER 65532:65532

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
ENV GOTRACEBACK=single

ENTRYPOINT ["/services/catalogs-read"]