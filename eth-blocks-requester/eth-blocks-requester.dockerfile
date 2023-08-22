FROM golang:1.19

WORKDIR /app/eth-blocks-requester

RUN go install github.com/cosmtrek/air@latest

COPY eth-helpers /app/eth-helpers

COPY eth-blocks-requester/go.mod /app/eth-blocks-requester
COPY eth-blocks-requester/.air.toml /app/eth-blocks-requester
COPY eth-blocks-requester/.env /app/eth-blocks-requester

RUN go mod download

COPY eth-blocks-requester /app/eth-blocks-requester

CMD ["air", "-c", ".air.toml"]