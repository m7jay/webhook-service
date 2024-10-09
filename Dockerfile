FROM golang:1.22-alpine3.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o main ./cmd

EXPOSE 8080

CMD ["./main"]