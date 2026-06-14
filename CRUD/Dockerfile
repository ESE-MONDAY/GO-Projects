# Step 1: Build the Go binary using the official Go image
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /movies-api main.go

# Step 2: Run the binary in a lightweight, clean environment
FROM alpine:latest
WORKDIR /
COPY --from=builder /movies-api /movies-api
EXPOSE 8080
CMD ["/movies-api"]