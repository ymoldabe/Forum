FROM golang:alpine AS build
ENV CGO_ENABLED=1

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

WORKDIR /app

COPY . .

RUN go build -o forum ./cmd/main.go
FROM alpine:latest
WORKDIR /app
COPY --from=build /app .
CMD ["./forum"]