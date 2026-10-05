FROM golang:1.27-bookworm AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/build/apid ./cmd/apid

FROM debian:bookworm-slim
WORKDIR /app

# need to use curl cli for health check
RUN apt-get update \
  && apt-get install -y --no-install-recommends curl ca-certificates \
  && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/build/apid ./apid

EXPOSE 5001 50051
CMD ["./apid"]
