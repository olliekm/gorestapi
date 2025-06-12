# ── BUILDER: force Go to cross-compile for linux/arm64 ───────
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# override Go’s target to exactly linux/arm64
ENV CGO_ENABLED=0 \
    GOOS=linux   \
    GOARCH=arm64

RUN go build -ldflags="-s -w" -o api     ./cmd/api
RUN go build -ldflags="-s -w" -o migrate ./cmd/migrate

# ── RUNTIME: minimal Alpine ─────────────────────────────────
FROM alpine:3.18
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/api     /usr/local/bin/api
COPY --from=builder /app/migrate /usr/local/bin/migrate
COPY --from=builder /app/cmd/migrate/migrations /migrations

RUN chmod +x /usr/local/bin/api /usr/local/bin/migrate

WORKDIR /
EXPOSE 8080
ENTRYPOINT ["api"]
