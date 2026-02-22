FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o app ./cmd/server

EXPOSE 8081
CMD ["./app"]