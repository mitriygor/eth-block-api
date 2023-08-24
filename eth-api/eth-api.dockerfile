FROM golang:1.19

WORKDIR /app/eth-api

RUN go install github.com/cosmtrek/air@latest

COPY eth-helpers /app/eth-helpers

COPY eth-api/go.mod /app/eth-api
COPY eth-api/.air.toml /app/eth-api

RUN go mod download

COPY eth-api /app/eth-api

EXPOSE 3000

CMD ["air", "-c", ".air.toml"]