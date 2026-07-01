FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/user-service /user-service

EXPOSE 8081

CMD ["/user-service"]
