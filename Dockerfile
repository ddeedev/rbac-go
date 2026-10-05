FROM golang:1.25-bookworm AS builder

WORKDIR /app

ADD . /app

RUN go build -o build/apid cmd/apid

FROM debian:bookworm-slim

COPY --from=builder app/build/apid .

EXPOSE 5001
EXPOSE 50051

CMD [ "./apid" ]
