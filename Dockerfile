# syntax=docker/dockerfile:1.7
FROM --platform=$BUILDPLATFORM golang:1.26.0-alpine3.22 AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

COPY . .
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags "-s -w" -o /out/service ./cmd/service

FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=builder /out/service /app/service

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/service"]
