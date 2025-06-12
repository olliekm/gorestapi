# ── STAGE 1: compile both tools ────────────────────────────
FROM golang:1.23-alpine AS builder
WORKDIR /app

# download deps
COPY go.mod go.sum Makefile ./
RUN apk add --no-cache make
RUN go mod download

# copy source & build
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=arm64
RUN make build                     # produces /app/bin/gorestapi
RUN go build -o bin/migrate cmd/migrate/main.go

# ── STAGE 2: runtime ───────────────────────────────────────
FROM alpine:3.18
# for nc, etc.
RUN apk add --no-cache ca-certificates netcat-openbsd
# use /app so that `file://cmd/migrate/migrations` resolves
WORKDIR /app

# bring in the API and migrate binaries
COPY --from=builder /app/bin/gorestapi /usr/local/bin/api
COPY --from=builder /app/bin/migrate   /usr/local/bin/migrate

# copy your SQL migrations into exactly cmd/migrate/migrations
COPY --from=builder /app/cmd/migrate/migrations ./cmd/migrate/migrations

RUN chmod +x /usr/local/bin/api /usr/local/bin/migrate

# default is to run the API; migrations will be kicked off by Compose
EXPOSE 8080
ENTRYPOINT ["api"]
