FROM golang:1.27-bookworm AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/build/apid ./cmd/apid

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/build/apid ./apid

EXPOSE 5001 50051
CMD ["./apid"]
