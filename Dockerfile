FROM golang:1.25.3-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app "./cmd"

FROM alpine:3.20

RUN apk --no-cache add ca-certificates curl && adduser -D -u 1000 appuser

WORKDIR /app

COPY --from=builder /build/app .

USER appuser

EXPOSE 8080 8081

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD curl -f http://localhost:8081/health || exit 1

ENTRYPOINT ["./app"]
