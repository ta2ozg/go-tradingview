# Build
FROM golang:1.24-alpine AS builder

WORKDIR /app


RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o trading-api .

FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/trading-api /trading-api
EXPOSE 3000

ENTRYPOINT ["/trading-api"]
