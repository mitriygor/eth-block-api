FROM golang:1.19

WORKDIR /app/eth-blocks-recorder

RUN go install github.com/cosmtrek/air@latest

COPY eth-helpers /app/eth-helpers

COPY eth-blocks-recorder/go.mod /app/eth-blocks-recorder
COPY eth-blocks-recorder/.air.toml /app/eth-blocks-recorder

RUN go mod download

COPY eth-blocks-recorder /app/eth-blocks-recorder

CMD ["air", "-c", ".air.toml"]