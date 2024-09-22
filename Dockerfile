# Build stage
FROM golang:1.20.6 AS builder

WORKDIR /go/src/tnyclick
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Final stage
FROM scratch

LABEL org.opencontainers.image.title="TnyClick"
LABEL org.opencontainers.image.description="A simple Go application for image uploading."
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.url="https://github.com/evanjarrett/tnyclick"
LABEL org.opencontainers.image.source="https://github.com/evanjarrett/tnyclick"

COPY --from=builder /go/src/tnyclick/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
