# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (required for go get)
RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o leaderboard-api main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/leaderboard-api .

COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./leaderboard-api"]
