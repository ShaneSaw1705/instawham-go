FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go mod tidy

COPY . .

RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]
