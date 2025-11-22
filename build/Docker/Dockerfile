FROM golang:1.24.0-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main cmd/telegram_app/main.go

CMD [ "./main" ]
