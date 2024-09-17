# Build the application
ARG GO_VERSION=1.23

FROM golang:${GO_VERSION} AS builder

ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o server ./cmd/server

# Create application container
FROM gcr.io/distroless/base-debian11 AS application

WORKDIR /
COPY --from=builder /src/server /server

USER nonroot:nonroot
ENTRYPOINT [ "/server" ]
