FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /map-api ./cmd/map-api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates \
 && addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /map-api .

USER appuser

EXPOSE 8090
ENTRYPOINT ["/app/map-api"]
