
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /movies-api main.go
FROM alpine:latest
WORKDIR /
COPY --from=builder /movies-api /movies-api
EXPOSE 8080
CMD ["/movies-api"]