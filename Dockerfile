FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o hyperfulcrum ./cmd/hyperfulcrum

FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /app/hyperfulcrum .
COPY --from=builder /app/migrations ./migrations  

EXPOSE 8080

CMD ["./hyperfulcrum"]
