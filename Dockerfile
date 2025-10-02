FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o comment-tree ./cmd/comment-tree/main.go

EXPOSE 8080