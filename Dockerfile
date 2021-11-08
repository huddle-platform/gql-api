FROM golang:1.16-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY graph ./graph
COPY server.go ./
RUN go build -o /app/main

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main /app/main
CMD ["/app/main"]