FROM golang:1.19
WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod ./
COPY .air.toml ./
RUN go mod download

COPY . .

EXPOSE 3000

CMD ["air", "-c", ".air.toml"]