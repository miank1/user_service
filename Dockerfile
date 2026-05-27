# -----------------------------
# Stage 1: Build user-service
# -----------------------------
FROM golang:1.24.6-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src

# Copy root module files
COPY go.mod go.sum ./

RUN go mod download

# Copy full monorepo source
COPY services ./services
COPY pkg ./pkg

# Build user-service binary
WORKDIR /src/services/user-service/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service main.go


# -----------------------------
# Stage 2: Runtime
# -----------------------------
FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /user-service /usr/local/bin/user-service

EXPOSE 8081

CMD ["/usr/local/bin/user-service"]