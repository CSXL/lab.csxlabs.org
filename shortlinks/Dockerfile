# Use base golang image from Docker Hub
FROM golang:1.20 AS build

WORKDIR /app

# Install dependencies in go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the application source code
COPY . ./

# Compile the application to /app.
# Skaffold passes in debug-oriented compiler flags
ARG SKAFFOLD_GO_GCFLAGS
RUN echo "Go gcflags: ${SKAFFOLD_GO_GCFLAGS}" \
    go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -mod=readonly -v -o /app

# Now create separate deployment image
FROM gcr.io/distroless/base

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single
ENV PORT=8080

# Copy template & assets
WORKDIR /shortlinks
COPY --from=build /app ./app
COPY templates templates/
COPY config.yaml config.yaml

ENTRYPOINT ["./app"]
