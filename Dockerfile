FROM golang:1.24.0-alpine3.21 AS builder

WORKDIR /app

COPY . .

RUN apk add --no-cache gcc musl-dev curl
ENV CGO_ENABLED=1 \
    GOCACHE=/go-cache \
    GOMODCACHE=/gomod-cache

RUN --mount=type=cache,target=/gomod-cache go mod download 

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache  \
    go build -ldflags="-s -w" -o main ./cmd/app

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config /app/config
COPY --from=builder /app/pkg/emails /app/pkg/emails

EXPOSE 8000

CMD ["/app/main"]