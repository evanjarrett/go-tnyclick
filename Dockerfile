# Build stage
FROM golang:1.20.6 AS builder

WORKDIR /go/src/tnyclick
COPY tnyclick .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Final stage
FROM scratch

COPY --from=builder /go/src/tnyclick/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]