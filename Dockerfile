FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/


FROM alpine:latest

WORKDIR /app

# copy binary
COPY --from=builder /app/app .

# copy .env
COPY .env .

EXPOSE 8081

CMD ["./app"]