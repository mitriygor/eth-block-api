FROM golang:1.19

WORKDIR /app/eth-redis-recorder

RUN go install github.com/cosmtrek/air@latest

COPY eth-helpers /app/eth-helpers

COPY eth-redis-recorder/go.mod /app/eth-redis-recorder
COPY eth-redis-recorder/.air.toml /app/eth-redis-recorder
COPY eth-redis-recorder/.env /app/eth-redis-recorder

RUN go mod download

COPY eth-redis-recorder /app/eth-redis-recorder

CMD ["air", "-c", ".air.toml"]