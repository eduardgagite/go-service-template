FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o service ./cmd/service

FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=builder /app/service /app/service

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/service"]
