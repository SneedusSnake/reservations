FROM golang:1.24.0-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o bin/telegram cmd/telegram_app/main.go

CMD [ "./bin/telegram" ]
