# Stage 1: build Go backend
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /map-api ./cmd/map-api

# Stage 2: run
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /map-api .
EXPOSE 8090
ENTRYPOINT ["/app/map-api"]
