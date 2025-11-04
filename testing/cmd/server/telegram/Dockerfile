FROM golang:1.24.0-alpine

WORKDIR /app

#COPY go.mod ./

#RUN go mod download

COPY . .

RUN go mod init github.com/SneedusSnake/telegram-fake-server

RUN go build -o svr ./main.go

CMD [ "./svr" ]
