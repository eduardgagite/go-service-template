# syntax=docker/dockerfile:1.7
FROM --platform=$BUILDPLATFORM golang:1.26.0-alpine3.22 AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go run github.com/swaggo/swag/cmd/swag@v1.16.5 init -g cmd/service/main.go -o docs --parseInternal

ARG TARGETOS TARGETARCH
ARG VERSION=dev COMMIT=none BUILD_DATE=unknown
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath \
    -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildDate=${BUILD_DATE}" \
    -o /out/service ./cmd/service

FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=builder /out/service /app/service
COPY --from=builder /src/docs /app/docs

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/service"]
