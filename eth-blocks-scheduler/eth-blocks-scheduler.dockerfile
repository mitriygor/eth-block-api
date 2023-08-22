FROM golang:1.19

WORKDIR /app/eth-blocks-scheduler

RUN go install github.com/cosmtrek/air@latest

COPY eth-helpers /app/eth-helpers

COPY eth-blocks-scheduler/go.mod /app/eth-blocks-scheduler
COPY eth-blocks-scheduler/.air.toml /app/eth-blocks-scheduler
COPY eth-blocks-scheduler/.env /app/eth-blocks-scheduler

RUN go mod download

COPY eth-blocks-scheduler /app/eth-blocks-scheduler

CMD ["air", "-c", ".air.toml"]