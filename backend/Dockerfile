FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /recly-api ./cmd/api

# Стадия 2: минимальный образ
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /recly-api /usr/local/bin/recly-api

EXPOSE 8080
ENTRYPOINT ["recly-api"]