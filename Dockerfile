# Build vault-raft-backup binary
FROM docker.io/library/golang:1.20 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

# Copy Go module manifests for dependency management
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy all project source files into the working directory
COPY . .

# Compile the application. The platform defaults are not specified, allowing Docker to automatically
# determine them based on the build platform (e.g., linux/arm64 for Apple M1)
RUN CGO_ENABLED=0 GOOS="${TARGETOS:-linux}" GOARCH="${TARGETARCH}" go build -a -o vault-raft-backup

# Switch to a minimal Distroless image to create the final image
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
ENTRYPOINT ["/vault-raft-backup"]

# Run the container as a non-root user for increased security
USER 65532:65532
WORKDIR /
COPY --from=builder /workspace/vault-raft-backup .
