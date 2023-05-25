FROM golang:1.20.4-alpine3.18

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

CMD ["./main"]