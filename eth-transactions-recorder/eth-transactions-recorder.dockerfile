FROM golang:1.19

WORKDIR /app/eth-transactions-recorder

RUN go install github.com/cosmtrek/air@latest

COPY eth-helpers /app/eth-helpers

COPY eth-transactions-recorder/go.mod /app/eth-transactions-recorder
COPY eth-transactions-recorder/.air.toml /app/eth-transactions-recorder

RUN go mod download

COPY eth-transactions-recorder /app/eth-transactions-recorder

CMD ["air", "-c", ".air.toml"]