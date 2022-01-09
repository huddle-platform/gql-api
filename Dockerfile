FROM golang:1.16-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY graph ./graph
COPY auth ./auth
COPY sql ./sql
COPY server.go ./
RUN go build -o /app/main

# Add certificate
# RUN apk --no-cache add curl
# RUN curl --create-dirs -o $HOME/.postgresql/root.crt -O https://cockroachlabs.cloud/clusters/741a8c79-a896-4f1f-963a-ac16126ff0bb/cert

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main /app/main
# COPY --from=builder /root/.postgresql/root.crt /root/.postgresql/root.crt
CMD ["/app/main"]