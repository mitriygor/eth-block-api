FROM golang:1.19

WORKDIR /app/eth-transactions-scheduler

RUN go install github.com/cosmtrek/air@latest

COPY eth-helpers /app/eth-helpers

COPY eth-transactions-scheduler/go.mod /app/eth-transactions-scheduler
COPY eth-transactions-scheduler/.air.toml /app/eth-transactions-scheduler

RUN go mod download

COPY eth-transactions-scheduler /app/eth-transactions-scheduler

CMD ["air", "-c", ".air.toml"]