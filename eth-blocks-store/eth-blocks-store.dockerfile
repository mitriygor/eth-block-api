FROM golang:1.19

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod ./
COPY .air.toml ./
COPY .env ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]