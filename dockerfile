FROM golang:1.23-alpine AS builder

WORKDIR /app


# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .


# build the API server
RUN go build -o api   ./cmd/api


# build the migration runner
RUN go build -o migrate ./cmd/migrate


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root
# pull in your two binaries
COPY --from=builder /app/api     .
COPY --from=builder /app/migrate .
# and your migrations so migrate can find them
COPY --from=builder /app/cmd/migrate/migrations ./cmd/migrate/migrations

RUN chmod +x ./api ./migrate
# expose your API port
EXPOSE 8080

# default entrypoint is your API server;
# we’ll override this in compose for migrations
ENTRYPOINT ["./api"]