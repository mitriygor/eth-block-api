FROM golang:1.19

WORKDIR /app/eth-transactions-requester

RUN go install github.com/cosmtrek/air@latest

COPY eth-helpers /app/eth-helpers

COPY eth-transactions-requester/go.mod /app/eth-transactions-requester
COPY eth-transactions-requester/.air.toml /app/eth-transactions-requester
COPY eth-transactions-requester/.env /app/eth-transactions-requester

RUN go mod download

COPY eth-transactions-requester /app/eth-transactions-requester

CMD ["air", "-c", ".air.toml"]